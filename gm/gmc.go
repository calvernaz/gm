package gm

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	
	"github.com/opentracing/opentracing-go/log"
	"go4.org/xdgdir"
)

var (
	configPath = filepath.Join("gm", "gm.json")
)

type GitManagerFile struct {
	Repositories []Repository `json:"repositories"`
}

type GitManagerConfig struct {
	Version string
	
	path string
	file *os.File
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
	
	gmc.file = file
	gmc.path = file.Name()
	
	return err
}

// Close close the configuration file
func (gmc *GitManagerConfig) Close() {
	if gmc.file != nil {
		err := gmc.file.Close()
		if err != nil {
			log.Error(err)
		}
	}
}
func (gmc *GitManagerConfig) File() *os.File {
	return gmc.file
}

func (gmc *GitManagerConfig) Path() string {
	return gmc.path
}