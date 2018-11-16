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

type GitManagerConfig struct {
	Version string
	
	path string
	file *os.File
}

// Open creates the config file if it doesn't exist, opens it otherwise.
// It returns an error in case it can create or open the file, the caller is
// responsible for close the file.
func (gmc *GitManagerConfig) Open(path string) error {
	if path == "" {
		path = configPath
	}
	
	f, err := xdgdir.Config.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			if file, err := xdgdir.Config.Create(path); err == nil {
				gmc.file = file
				gmc.path = path
				return nil
			}
			return errors.New("failed to create config file")
		}
		return err
	}
	gmc.file = f
	gmc.path = path
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
