package manager

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"


	"go4.org/xdgdir"
)

const (
	Succeed = "updated"
	Failed = "not updated"
)

var (
	configPath = filepath.Join("gm", "gm.json")
)

type GitManagerFile struct {
	Repositories []RepositoryEntry `json:"repositories"`
}

type GitManagerConfig struct {
	Version string

	file *os.File

	wg sync.WaitGroup
	ch chan Operation

}

// Open creates the config file if it doesn't exist, opens it otherwise.
// It returns an error in case it can create or open the file, the caller is
// responsible for close the file.
func (gmc *GitManagerConfig) Open(path string) (err error) {
	if path == "" {
		path = configPath
	}

	var file *os.File
	file, err = xdgdir.Config.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = xdgdir.Config.Create(path)
			if err != nil {
				return errors.New("failed to create config file")
			}
		}

		if err != nil {
			return errors.New("failed to open the config file")
		}
	}

	// configuration file
	gmc.file = file

	return err
}

// Close close the configuration file
func (gmc *GitManagerConfig) Close() {
	if gmc.file != nil {
		err := gmc.file.Close()
		if err != nil {
			Error.Printf("failed to close file: %v", err)
		}
	}

	close(gmc.ch)
}

// File the configuration file description
func (gmc *GitManagerConfig) File() *os.File {
	return gmc.file
}

// Loop loops waiting for operations to execute
func (gmc *GitManagerConfig) Loop() {
	gmc.ch = make(chan Operation, 2)
	for {
		select {
		case op, ok := <-gmc.ch:
			if ok {
				err := op.Execute()
				if err != nil {
					//gmc.bl.Buffer(op.Repo.Name, Failed, err.Error())
					ErrorLog("failed to update repository: %v", op.Repo.Name)
				} else {
					//gmc.bl.Buffer(op.Repo.Name, Succeed)
					InfoLog("updated repository: %v", op.Repo.Name)
				}
			}
			gmc.wg.Done()
		}
	}
}

// Add new operation to run
func (gmc *GitManagerConfig) Run(operation Operation) {
	gmc.wg.Add(1)
	gmc.ch <- operation
}

func (gmc *GitManagerConfig) Wait() {
	gmc.wg.Wait()
}

func (gmc *GitManagerConfig) RemoveDups(entries []RepositoryEntry) []RepositoryEntry {
	bytes, _ := ioutil.ReadAll(gmc.file)

	var configFile GitManagerFile
	_ = json.Unmarshal(bytes, &configFile)

	all := append(configFile.Repositories, entries...)
	j := 0
	for i := 1; i < len(all); i++ {
		if filepath.Clean(all[j].Path) == filepath.Clean(all[i].Path) {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		all[j] = all[i]
	}
	result := all[:j+1]
	return result
}
