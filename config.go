package gm

import (
	"os"

	"github.com/opentracing/opentracing-go/log"
	"go4.org/xdgdir"
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
	if f, err := xdgdir.Config.Create(path); err == nil {
		gmc.path = path
		gmc.file = f
	} else {
		return err
	}
	return nil
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
