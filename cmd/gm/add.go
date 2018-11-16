package main

import (
	"flag"
	"github.com/calvernaz/gm/gm"
	"github.com/calvernaz/gm/subcmd"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (s *State) add(args ...string) {
	const help = `
`
	
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	verbose := fs.Bool("v", false, "log each repository as it is added")
	recur := fs.Bool("R", false, "recursively add repositories")
	overwrite := fs.Bool("overwrite", true, "overwrite existing repositories")
	
	s.ParseFlags(fs, args, help, "add repository...")
	cs := &copyState{
		state:     s,
		flagSet:   fs,
		overwrite: *overwrite,
		recur:     *recur,
		verbose:   *verbose,
	}
	
	var repositories []gm.Repository
	for _, repo := range fs.Args() {
		repositories = append(repositories, cs.glob(repo)...)
	}
	
	s.copyCommand(cs, repositories)
}

type copyState struct {
	state     *State
	flagSet   *flag.FlagSet // Used only to call Usage.
	overwrite bool
	recur     bool
	verbose   bool
}

func (s *State) copyCommand(cs *copyState, repositories []gm.Repository) {

	if len(repositories) != 1 {
		s.Exitf("copying multiple files but %s is not a directory", dstFile.path)
		usageAndExit(cs.flagSet)
	}
	
	s.copyToDir(cs, repositories)
}

// copyToDir copies the source files to the destination directory.
// It recurs if -R is set and a source is a subdirectory.
func (s *State) copyToDir(cs *copyState, src []gm.Repository) {
	for _, from := range src {
	
	}
}

// isDir reports whether the file is a directory either in Upspin
// or in the local file system.
func (s *State) isDir(cf gm.Repository) bool {
	info, err := os.Stat(cf.Path)
	return err == nil && info.IsDir()
}

// glob glob-expands the argument, which could be a local file
// name. Files on the local machine
// must be identified by absolute paths.
// That is, they must be full paths.
func (cs *copyState) glob(pattern string) (files []gm.Repository) {
	if pattern == "" {
		cs.state.Exitf("empty path name")
	}
	
	// Path on local machine?
	if isLocal(pattern) {
		for _, path := range cs.state.GlobLocal(subcmd.Tilde(pattern)) {
			files = append(files, gm.Repository{
				Name:       "",
				Enabled:    false,
				LastUpdate: time.Time{},
				Path: path,
			})
		}
		return files
		
	}
	
	// Extra check to catch use of relative path on local machine.
	if !strings.Contains(pattern, "@") {
		cs.state.Exitf("local pattern not qualified path: %s", pattern)
	}

	return files
}

// isLocal reports whether the argument names a fully-qualified local file.
// TODO: This is Unix-specific.
func isLocal(file string) bool {
	switch {
	case filepath.IsAbs(file):
		return true
	case file == ".", file == "..":
		return true
	case strings.HasPrefix(file, "~"):
		return true
	case strings.HasPrefix(file, "./"):
		return true
	case strings.HasPrefix(file, "../"):
		return true
	}
	return false
}
