package main

import (
	"flag"
	"os"
	"time"
	
	"github.com/calvernaz/gm/manager"
)

func (s *State) get(args ...string) {
	const help = ``
	
	// commands
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	s.ParseFlags(fs, args, help, "get [repository]")
	
	wd, _ := os.Getwd()
	
	for _, repo := range args {
		repo := manager.RepositoryEntry{
			Name:       repo,
			Enabled:    true,
			LastUpdate: time.Now(),
			Path:       wd,
		}
		op := manager.Operation{
			Repo:   repo,
			OpType: manager.Download,
		}
		s.gmc.Run(op)
	}
}
