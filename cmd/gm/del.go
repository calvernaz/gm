package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/calvernaz/gm/manager"
)

type deleteState struct {
	state     *State
	flagSet   *flag.FlagSet // Used only to call Usage.
}

func (s *State) del(args ...string) {
	const help = ``

	// commands
	fs := flag.NewFlagSet("del", flag.ExitOnError)

	cs := &deleteState{
		state:     s,
		flagSet:   fs,
	}
	s.ParseFlags(fs, args, help, "delete repository")
	s.deleteCommand(cs, fs.Args())
}

func (s *State) deleteCommand(cs *deleteState, repositoryNames []string) {
	f := s.gmc.File()
	_, err := f.Stat()
	if err != nil {
		s.Exitf("failed to obtain file info")
		usageAndExit(cs.flagSet)
	}

	content, err := ioutil.ReadFile(s.gmc.File().Name())
	if err != nil {
		s.Exitf("failed reading the configuration file: %v", err)
	}

	var repos manager.GitManagerFile
	err = json.Unmarshal(content, &repos)
	if err != nil {
		s.Exitf("failed to read repositories", err)
	}

	for i, repo := range repos.Repositories {
		if ContainsString(repositoryNames, repo.Name) {
			repos.Repositories = append(repos.Repositories[:i], repos.Repositories[i+1:]...)
		}
	}

	b, err := json.Marshal(repos)
	if err == nil {
		err = ioutil.WriteFile(s.gmc.File().Name(), b, 0644)
		if err != nil {
			s.Exitf("failed to add repository: %v", err)
		}
		for _, repositoryName := range repositoryNames {
			fmt.Printf("%q deleted", repositoryName)
		}
		return
	}
}

func ContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
