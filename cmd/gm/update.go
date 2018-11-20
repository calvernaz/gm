package main

import (
	"encoding/json"
	"fmt"
	"github.com/calvernaz/gm/gm"
	"github.com/calvernaz/gm/subcmd"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"strings"
)

func (s *State) update(args ...string) {
	const help = ``
	
	content, err := ioutil.ReadFile(s.gmc.Path())
	if err != nil {
		s.Exitf("failed reading the configuration file: %v", err)
	}
	
	repos := gm.GitManagerFile{}
	err = json.Unmarshal(content, &repos)
	if err != nil {
		s.Exitf("failed to read repositories", err)
	}
	
	Info("git pull origin")
	for _, repo := range repos.Repositories {
		r, err := git.PlainOpen(subcmd.Tilde(repo.Path))
		CheckIfError(err)
		w, err := r.Worktree()
		CheckIfError(err)
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		ref, err := r.Head()
		commit, err := r.CommitObject(ref.Hash())
		CheckIfError(err)
		
		fmt.Println(commit)
	}
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
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
