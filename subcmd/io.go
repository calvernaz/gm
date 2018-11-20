package subcmd

import (
	"github.com/pkg/errors"
	"os"
	osUser "os/user"
	"path/filepath"
	"strings"
)

var userLookup = osUser.Lookup
var home string // Main user's home directory.

// Homedir returns the home directory of the OS' logged-in user.
// TODO(adg): move to osutil package?
func Homedir() (string, error) {
	u, err := osUser.Current()
	// user.Current may return an error, but we should only handle it if it
	// returns a nil user. This is because os/user is wonky without cgo,
	// but it should work well enough for our purposes.
	if u == nil {
		e := errors.New("lookup of current user failed")
		if err != nil {
			e = errors.Errorf("%v: %v", e, err)
		}
		return "", e
	}
	h := u.HomeDir
	if h == "" {
		return "", errors.New("user home directory not found")
	}
	if err := isDir(h); err != nil {
		return "", err
	}
	return h, nil
}

func isDir(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		return errors.WithMessage(err, "External I/O error such as network failure.")
	}
	if !fi.IsDir() {
		return errors.Errorf("Item does not exist: %v", p)
	}
	return nil
}

func homeDir(who string) string {
	if who == "" {
		if home == "" {
			var err error
			home, err = Homedir()
			if err != nil {
				return "~" // What else can we do?
			}
			return home
		}
	}
	u, err := userLookup(who)
	if err != nil {
		return "~" + who // Again, what else can we do?
	}
	return u.HomeDir
}
// Tilde processes a leading tilde, if any, in the local file name.
// If the file name does not begin with a tilde, Tilde returns the argument unchanged.
// This special processing (only) is applied to all local file names passed to
// functions in this package.
// If the target user does not exist, it returns the original string.
func Tilde(file string) string {
	if file == "" || file[0] != '~' {
		return file
	}
	if file == "~" {
		return homeDir("")
	}
	slash := strings.IndexByte(file, '/')
	if slash < 0 {
		return homeDir(file[1:])
	}
	return filepath.Join(homeDir(file[1:slash]), file[slash+1:])
}


// HasGlobChar reports whether the string contains an unescaped Glob metacharacter.
func HasGlobChar(pattern string) bool {
	esc := false
	for _, r := range pattern {
		if esc {
			esc = false // TODO: What if next rune is '/'?
			continue
		}
		switch r {
		case '\\':
			esc = true
		case '*', '[', '?':
			return true
		}
	}
	return false
}

// GlobLocal glob-expands the argument, which should be a syntactically
// valid Glob pattern (including a plain file name). If the pattern is erroneous
// or matches no files, the function exits.
func (s *State) GlobLocal(pattern string) []string {
	pattern = Tilde(pattern)
	// If it has no metacharacters, leave it alone.
	if !HasGlobChar(pattern) {
		return []string{pattern}
	}
	strs, err := filepath.Glob(pattern)
	if err != nil {
		// Bad pattern.
		s.Exitf("bad local Glob pattern %q: %v", pattern, err)
	}
	if len(strs) == 0 {
		s.Exitf("no path matches %q", pattern)
	}
	return strs
}
