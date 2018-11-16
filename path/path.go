package path

import (
	"strings"
)

// Parse parses a full file name, including the user, validates it,
// and returns its parsed form. If the name is a user root directory,
// the trailing slash is optional. The name is 'cleaned' (see the Clean
// function) to canonicalize it.
func Parse(pathName upspin.PathName) (Parsed, error) {
	name := string(pathName)
	// Pull off the user name.
	var userName string
	slash := strings.IndexByte(name, '/')
	if slash < 0 {
		userName = name
	} else {
		userName = name[:slash]
	}
	if _, _, _, err := user.Parse(upspin.UserName(userName)); err != nil {
		// Bad user name.
		return Parsed{}, err
	}
	p := Parsed{
		// If pathName is already clean, which it usually is, this will not allocate.
		path: Clean(pathName),
	}
	return p, nil
}
