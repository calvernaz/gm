package gm

import (
	"time"
	"github.com/calvernaz/gm/subcmd"
	"fmt"
	"os"
	"gopkg.in/src-d/go-git.v4"
)

// Repository format entry
type Repository struct {
	// Repository Name
	Name string `json:"name"`
	// Tracking repository
	Enabled bool `json:"enabled"`
	// Last repository update
	LastUpdate time.Time `json:"last_update"`
	// Git Repository path
	Path string
}

func (r Repository) Update() {
	r, err := git.PlainOpen(subcmd.Tilde(r.Path))
	CheckIfError(err)
	w, err := r.Worktree()
	CheckIfError(err)
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	ref, err := r.Head()
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

