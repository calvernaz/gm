package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type vcsCmd struct {
	name string
	cmd  string // name of binary to invoke command

	createCmd   []string // commands to download a fresh copy of a repository
	downloadCmd []string // commands to download updates into an existing repository

	tagCmd         []tagCmd // commands to list tags
	tagLookupCmd   []tagCmd // commands to lookup tags before running tagSyncCmd
	tagSyncCmd     []string // commands to sync to specific tag
	tagSyncDefault []string // commands to sync to default tag

	scheme  []string
	pingCmd string

	remoteRepo  func(v *vcsCmd, rootDir string) (remoteRepo string, err error)
	resolveRepo func(v *vcsCmd, rootDir, remoteRepo string) (realRepo string, err error)
}

var defaultSecureScheme = map[string]bool{
	"https":   true,
	"git+ssh": true,
	"bzr+ssh": true,
	"svn+ssh": true,
	"ssh":     true,
}

func (v *vcsCmd) isSecure(repo string) bool {
	u, err := url.Parse(repo)
	if err != nil {
		// If repo is not a URL, it's not secure.
		return false
	}
	return v.isSecureScheme(u.Scheme)
}

func (v *vcsCmd) isSecureScheme(scheme string) bool {
	switch v.cmd {
	case "git":
		// GIT_ALLOW_PROTOCOL is an environment variable defined by Git. It is a
		// colon-separated list of schemes that are allowed to be used with git
		// fetch/clone. Any scheme not mentioned will be considered insecure.
		if allow := os.Getenv("GIT_ALLOW_PROTOCOL"); allow != "" {
			for _, s := range strings.Split(allow, ":") {
				if s == scheme {
					return true
				}
			}
			return false
		}
	}
	return defaultSecureScheme[scheme]
}

// A tagCmd describes a command to list available tags
// that can be passed to tagsSyncCmd
type tagCmd struct {
	cmd     string // command to list tags
	pattern string // regexp to extract tags from list
}

// vcsList lists the known version control systems
var vcsList = []*vcsCmd{
	vcsGit,
}

// VcsByCmd returns the version control system for the given
// command name (hg, git, svn, bzr).
func VcsByCmd(cmd string) *vcsCmd {
	for _, vcs := range vcsList {
		if vcs.cmd == cmd {
			return vcs
		}
	}
	return nil
}

// vcsGit describes how to use Git.
var vcsGit = &vcsCmd{
	name: "Git",
	cmd:  "git",

	createCmd:   []string{"clone {repo} {dir}"},
	downloadCmd: []string{"pull --ff-only", "submodule update --init --recursive"},

	tagCmd: []tagCmd{
		// tags/xxx matches a git tag named xxx
		// origin/xxx matches a git branch named xxx on the default remote repository
		{"show-ref", `(?:tags|origin)/(\S+)$`},
	},
	tagLookupCmd: []tagCmd{
		{"show-ref tags/{tag} origin/{tag}", `((?:tags|origin)/\S+)$`},
	},
	tagSyncCmd: []string{"checkout {tag}", "submodule update --init --recursive"},
	// both createCmd and downloadCmd update the working dir.
	// No need to do more here. We used to 'checkout master'
	// but that doesn't work if the default branch is not named master.
	// DO NOT add 'checkout master' here.
	// See golang.org/issue/9032.
	tagSyncDefault: []string{"submodule update --init --recursive"},

	scheme:     []string{"git", "https", "http", "git+ssh", "ssh"},
	pingCmd:    "ls-remote {scheme}://{repo}",
	remoteRepo: gitRemoteRepo,
}

// scpSyntaxRe matches the SCP-like addresses used by Git to access
// repositories by SSH.
var scpSyntaxRe = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

func gitRemoteRepo(vcsGit *vcsCmd, rootDir string) (remoteRepo string, err error) {
	cmd := "config remote.origin.url"
	errParse := errors.New("unable to parse output of git " + cmd)
	errRemoteOriginNotFound := errors.New("remote origin not found")
	outb, err := vcsGit.run1(rootDir, cmd, nil, false)
	if err != nil {
		// if it doesn't output any message, it means the config argument is correct,
		// but the config value itself doesn't exist
		if outb != nil && len(outb) == 0 {
			return "", errRemoteOriginNotFound
		}
		return "", err
	}
	out := strings.TrimSpace(string(outb))

	var repoURL *url.URL
	if m := scpSyntaxRe.FindStringSubmatch(out); m != nil {
		// Match SCP-like syntax and convert it to a URL.
		// Eg, "git@github.com:user/repo" becomes
		// "ssh://git@github.com/user/repo".
		repoURL = &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}
	} else {
		repoURL, err = url.Parse(out)
		if err != nil {
			return "", err
		}
	}

	// Iterate over insecure schemes too, because this function simply
	// reports the state of the repo. If we can't see insecure schemes then
	// we can't report the actual repo URL.
	for _, s := range vcsGit.scheme {
		if repoURL.Scheme == s {
			return repoURL.String(), nil
		}
	}
	return "", errParse
}

func (v *vcsCmd) String() string {
	return v.name
}

// run runs the command line cmd in the given directory.
// keyval is a list of key, value pairs. run expands
// instances of {key} in cmd into value, but only after
// splitting cmd into individual arguments.
// If an error occurs, run prints the command line and the
// command's combined stdout+stderr to standard error.
// Otherwise run discards the command's output.
func (v *vcsCmd) run(dir string, cmd string, keyval ...string) error {
	_, err := v.run1(dir, cmd, keyval, true)
	return err
}

// runVerboseOnly is like run but only generates error output to standard error in verbose mode.
func (v *vcsCmd) runVerboseOnly(dir string, cmd string, keyval ...string) error {
	_, err := v.run1(dir, cmd, keyval, false)
	return err
}

// runOutput is like run but returns the output of the command.
func (v *vcsCmd) runOutput(dir string, cmd string, keyval ...string) ([]byte, error) {
	return v.run1(dir, cmd, keyval, true)
}

func (v *vcsCmd) run1(dir string, cmdline string, keyval []string, verbose bool) ([]byte, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}
	args := strings.Fields(cmdline)
	for i, arg := range args {
		args[i] = expand(m, arg)
	}

	_, err := exec.LookPath(v.cmd)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "missing %s command.", v.name)
		return nil, err
	}

	cmd := exec.Command(v.cmd, args...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		if verbose {
			_, _ = fmt.Fprintf(os.Stderr, "# cd %s; %s %s\n", dir, v.cmd, strings.Join(args, " "))
			if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
				_, _ = os.Stderr.Write(ee.Stderr)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, err.Error())
			}
		}
	}
	return out, err
}

// ping pings to determine scheme to use.
func (v *vcsCmd) ping(scheme, repo string) error {
	return v.runVerboseOnly(".", v.pingCmd, "scheme", scheme, "repo", repo)
}

// Create creates a new copy of repo in dir.
// The parent of dir must exist; dir must not.
func (v *vcsCmd) Create(dir, repo string) error {
	for _, cmd := range v.createCmd {
		if err := v.run(".", cmd, "dir", dir, "repo", repo); err != nil {
			return err
		}
	}
	return nil
}

// download downloads any new changes for the repo in dir.
func (v *vcsCmd) Download(dir string) error {
	for _, cmd := range v.downloadCmd {
		if err := v.run(dir, cmd); err != nil {
			return err
		}
	}
	return nil
}

// tags returns the list of available tags for the repo in dir.
func (v *vcsCmd) tags(dir string) ([]string, error) {
	var tags []string
	for _, tc := range v.tagCmd {
		out, err := v.runOutput(dir, tc.cmd)
		if err != nil {
			return nil, err
		}
		re := regexp.MustCompile(`(?m-s)` + tc.pattern)
		for _, m := range re.FindAllStringSubmatch(string(out), -1) {
			tags = append(tags, m[1])
		}
	}
	return tags, nil
}

// tagSync syncs the repo in dir to the named tag,
// which either is a tag returned by tags or is v.tagDefault.
func (v *vcsCmd) tagSync(dir, tag string) error {
	if v.tagSyncCmd == nil {
		return nil
	}
	if tag != "" {
		for _, tc := range v.tagLookupCmd {
			out, err := v.runOutput(dir, tc.cmd, "tag", tag)
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`(?m-s)` + tc.pattern)
			m := re.FindStringSubmatch(string(out))
			if len(m) > 1 {
				tag = m[1]
				break
			}
		}
	}

	if tag == "" && v.tagSyncDefault != nil {
		for _, cmd := range v.tagSyncDefault {
			if err := v.run(dir, cmd); err != nil {
				return err
			}
		}
		return nil
	}

	for _, cmd := range v.tagSyncCmd {
		if err := v.run(dir, cmd, "tag", tag); err != nil {
			return err
		}
	}
	return nil
}

// vcsFromDir inspects dir and its parents to determine the
// version control system and code repository to use.
// On return, root is the import path
// corresponding to the root of the repository.
func VcsFromDir(dir string) (vcs *vcsCmd, err error) {
	// Clean and double-check that dir is in (a subdirectory of) srcRoot.
	dir = filepath.Clean(dir)
	var vcsRet *vcsCmd
	var rootRet string

	origDir := dir

	for _, vcs := range vcsList {
		if _, err := os.Stat(filepath.Join(dir, "."+vcs.cmd)); err == nil {
			root := filepath.ToSlash(dir)
			// Record first VCS we find, but keep looking,
			// to detect mistakes like one kind of VCS inside another.
			if vcsRet == nil {
				vcsRet = vcs
				rootRet = root
				continue
			}
			// Allow .git inside .git, which can arise due to submodules.
			if vcsRet == vcs && vcs.cmd == "git" {
				continue
			}
			// Otherwise, we have one VCS inside a different VCS.
			return nil, fmt.Errorf("directory %q uses %s, but parent %q uses %s",
				filepath.Join(dir, rootRet), vcsRet.cmd, filepath.Join(dir, root), vcs.cmd)
		}
	}

	if vcsRet != nil {
		return vcsRet, nil
	}

	return nil, fmt.Errorf("directory %q is not using a known version control system", origDir)
}

// RepoRoot describes the repository root for a tree of source code.
type RepoRoot struct {
	Repo     string // repository URL, including scheme
	Root     string // import path corresponding to root of repo
	IsCustom bool   // defined by served <meta> tags (as opposed to hard-coded pattern)
	VCS      string // vcs type ("mod", "git", ...)

	vcs *vcsCmd // internal: vcs command access
}

var httpPrefixRE = regexp.MustCompile(`^https?:`)

// validateRepoRoot returns an error if repoRoot does not seem to be
// a valid URL with scheme.
func validateRepoRoot(repoRoot string) error {
	url, err := url.Parse(repoRoot)
	if err != nil {
		return err
	}
	if url.Scheme == "" {
		return errors.New("no scheme")
	}
	return nil
}

// expand rewrites s to replace {k} with match[k] for each key k in match.
func expand(match map[string]string, s string) string {
	// We want to replace each match exactly once, and the result of expansion
	// must not depend on the iteration order through the map.
	// A strings.Replacer has exactly the properties we're looking for.
	oldNew := make([]string, 0, 2*len(match))
	for k, v := range match {
		oldNew = append(oldNew, "{"+k+"}", v)
	}
	return strings.NewReplacer(oldNew...).Replace(s)
}

// noVCSSuffix checks that the repository name does not
// end in .foo for any version control system foo.
// The usual culprit is ".git".
func noVCSSuffix(match map[string]string) error {
	repo := match["repo"]
	for _, vcs := range vcsList {
		if strings.HasSuffix(repo, "."+vcs.cmd) {
			return fmt.Errorf("invalid version control suffix in %s path", match["prefix"])
		}
	}
	return nil
}

// bitbucketVCS determines the version control system for a
// Bitbucket repository, by using the Bitbucket API.
func bitbucketVCS(match map[string]string) error {
	if err := noVCSSuffix(match); err != nil {
		return err
	}

	var resp struct {
		SCM string `json:"scm"`
	}
	url := expand(match, "https://api.bitbucket.org/2.0/repositories/{bitname}?fields=scm")
	data, err := Get(url)
	if err != nil {
		if httpErr, ok := err.(*HTTPError); ok && httpErr.StatusCode == 403 {
			// this may be a private repository. If so, attempt to determine which
			// VCS it uses. See issue 5375.
			root := match["root"]
			for _, vcs := range []string{"git", "hg"} {
				if VcsByCmd(vcs).ping("https", root) == nil {
					resp.SCM = vcs
					break
				}
			}
		}

		if resp.SCM == "" {
			return err
		}
	} else {
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("decoding %s: %v", url, err)
		}
	}

	if VcsByCmd(resp.SCM) != nil {
		match["vcs"] = resp.SCM
		if resp.SCM == "git" {
			match["repo"] += ".git"
		}
		return nil
	}

	return fmt.Errorf("unable to detect version control system for bitbucket.org/ path")
}
