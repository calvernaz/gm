package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/calvernaz/gm/manager"
)

func (s *State) update(args ...string) {

	// reads the manager file
	content, err := ioutil.ReadFile(s.gmc.File().Name())
	if err != nil {
		s.Exitf("failed reading the configuration mf: %v", err)
	}

	// unmarshal json file
	mf := manager.GitManagerFile{}
	err = json.Unmarshal(content, &mf)
	if err != nil {
		s.Exitf("failed to read repositories", err)
	}

	// run update on all enabled repositories
	for _, repo := range mf.Repositories {
		if repo.Enabled {
			op := manager.Operation{
				Repo:   repo,
				OpType: manager.Update,
			}
			s.gmc.Run(op)
		}
	}
}
