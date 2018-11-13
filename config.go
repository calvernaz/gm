package gm

import (
	"os"
)

type GitManagerConfig struct {
	Version string
	
	path string
	file *os.File
}

// Open creates and opens a database at the given path.
// If the file does not exist then it will be created automatically.
// Passing in nil options will open the config with the default options.
func Open(path string, mode os.FileMode) (*GitManagerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	gm := &GitManagerConfig{
		"1",
		path,
		file,
	}
	return gm, nil
}
