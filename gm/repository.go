package gm

import (
	"fmt"
	"github.com/calvernaz/gm/subcmd"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"time"
	"path"
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

func (r Repository) Update() error {
	Info("updating repository %v", path.Base(r.Path))
	gitRepository, err := git.PlainOpen(subcmd.Tilde(r.Path))
	if err != nil {
		return err
	}
	//CheckIfError(err)
	w, err := gitRepository.Worktree()
	if err != nil {
		return err
	}
	
	//CheckIfError(err)
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}
	
	ref, err := gitRepository.Head()
	if err != nil {
		return err
	}
	
	commit, err := gitRepository.CommitObject(ref.Hash())
	if err != nil {
		return err
	}
	// CheckIfError(err)
	fmt.Println(commit)
	return nil
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

