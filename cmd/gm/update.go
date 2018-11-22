package main

import (
	"encoding/json"
	"github.com/calvernaz/gm/gm"
	"io/ioutil"
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
	
	for _, repo := range repos.Repositories {
//		go func(r gm.Repository) {
			s.gmc.Run(gm.Operation{
				Repo:   repo,
				OpType: gm.Update,
			})
//		}(repo)
	}
}
