package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/calvernaz/gm/manager"
)

func (s *State) get(args ...string) {
	const help = ``

	// commands
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	s.ParseFlags(fs, args, help, "get [repository] [path]")

	wd, _ := os.Getwd()
	if len(args) > 0 {
		wd = args[1]
	}

	fmt.Println(fs.Args())
	repo := manager.RepositoryEntry{
		Name:       args[0],
		Enabled:    true,
		LastUpdate: time.Now(),
		Path:       wd,
	}
	op := manager.Operation{
		Repo:   repo,
		OpType: manager.Download,
	}
	s.gmc.Run(op)

	cs := &copyState{
		state:     s,
		flagSet:   fs,
	}

	repoName := path.Base(repo.Name)
	repoName = repoName[0:len(repoName)-len(filepath.Ext(".git"))]

	var repositories []manager.RepositoryEntry
	repositories = append(repositories, cs.glob(filepath.Join(wd, repoName))...)

	s.copyCommand(cs, repositories)
}
