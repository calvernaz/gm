package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Go = "go"
)

var (
	quiet = flag.Bool("quiet", false, "Don't print anything unless there's a failure.")
	static = flag.Bool("static", false, "Build a static binary, so it can in a empty container.")
	stampVersion = flag.Bool("stampversion", true, "Stamp version into buildinfo.GitInfo.")

	gitVersionRx = regexp.MustCompile(`\b\d\d\d\d-\d\d-\d\d-[0-9a-f]{10}\b`)
	pkRoot string
)

func main()  {
	log.SetFlags(0)
	// parse arguments
	flag.Parse()
	// go install :)
	baseArgs := []string{"install", "-v"}
	// target
	target := []string {
		"github.com/calvernaz/gm/cmd/gm",
	}
	var ldFlags string
	// static link
	if *static {
		// -w omits dwarf table
		// -d disable generation of dynamic executables
		// linkmod "internal" - https://golang.org/src/cmd/cgo/doc.go
		ldFlags = "-w -d -linkmode internal"
	}
	// add version info
	if *stampVersion {
		if ldFlags != "" {
			ldFlags += " "
		}
		version := getVersion()
		gitRev := getGitVersion()
		ldFlags += "-X \"github.com/calvernaz/gm/buildinfo.GitInfo=" + gitRev + "\""
		ldFlags += "-X \"github.com/calvernaz/gm/buildinfo.Version=" + version + "\""
	}
	// if used before, append to base arguments
	if ldFlags != "" {
		baseArgs = append(baseArgs, "--ldflags="+ldFlags)
	}
	// add target arguments to base arguments
	args := append(baseArgs, target...)
	// setup go command
	cmd := exec.Command(Go, args...)

	// setup output
	var output bytes.Buffer
	if *quiet {
		cmd.Stdout = &output
		cmd.Stderr = &output
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// run the command
	if err := cmd.Run(); err != nil {
		log.Fatal("error building main binaries: %v\n%s", err, output.String())
	}
}

func getVersion() string {
	if _, err := exec.LookPath("git"); err != nil {
		return ""
	}
	if _, err := os.Stat(filepath.Join(pkRoot, ".git")); os.IsNotExist(err) {
		return ""
	}
	cmd := exec.Command("git", "describe", "--abbrev=0")
	cmd.Dir = pkRoot
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running git describe in %s: %v, tag isn't reachable from the commit", pkRoot, err)
	}
	v := strings.TrimSpace(string(out))
	return v
}

func getGitVersion() string {
	if _, err := exec.LookPath("git"); err != nil {
		return ""
	}
	if _, err := os.Stat(filepath.Join(pkRoot, ".git")); os.IsNotExist(err) {
		return ""
	}
	cmd := exec.Command("git", "rev-list", "--max-count=1", "--pretty=format:'%ad-%h'",
		"--date=short", "--abbrev=10", "HEAD")
	cmd.Dir = pkRoot
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running git rev-list in %s: %v", pkRoot, err)
	}
	v := strings.TrimSpace(string(out))
	if m := gitVersionRx.FindStringSubmatch(v); m != nil {
		v = m[0]
	} else {
		panic("Failed to find git version in " + v)
	}
	cmd = exec.Command("git", "diff", "--exit-code")
	cmd.Dir = pkRoot
	if err := cmd.Run(); err != nil {
		v += "+"
	}
	return v
}
