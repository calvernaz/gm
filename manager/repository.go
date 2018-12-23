package manager

import (
	"fmt"
	"os"
	"path"
	"time"
	
	"github.com/calvernaz/gm/internal"
	"github.com/calvernaz/gm/subcmd"
)

// RepositoryEntry format entry
type RepositoryEntry struct {
	// RepositoryEntry Name
	Name string `json:"name"`
	// Tracking repository
	Enabled bool `json:"enabled"`
	// Last repository update
	LastUpdate time.Time `json:"last_update"`
	// Git RepositoryEntry path
	Path string `json:"path"`
}

// Update updates the repository
func (r RepositoryEntry) Update() error {
	info("updating repository: %v", path.Base(r.Path))
	
	vcs, err := internal.VcsFromDir(subcmd.Tilde(r.Path))
	if err != nil {
		return err
	}
	//CheckIfError(err)
	
	err = vcs.Download(r.Path)
	if err != nil {
		return err
	}
	//CheckIfError(err)
	return nil
}

// info should be used to describe the example commands that are about to run.
func info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
