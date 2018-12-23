package manager

import (
	"os"
	"path/filepath"
	"sync"
	
	"github.com/calvernaz/gm/log"
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
	
	bl *log.BufferLog
	
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
	
	// buffered log
	gmc.bl = log.NewBufferLog()
	
	return err
}

// Close close the configuration file
func (gmc *GitManagerConfig) Close() {
	if gmc.file != nil {
		err := gmc.file.Close()
		if err != nil {
			log.Error.Printf("failed to close file: %v", err)
		}
	}

	close(gmc.ch)
}

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
					gmc.bl.Buffer(op.Repo.Name, Failed)
				} else {
					gmc.bl.Buffer(op.Repo.Name, Succeed)
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
	gmc.bl.Print()
}
