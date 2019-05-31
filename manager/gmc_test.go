package manager

import (
	"os"
	"path/filepath"
	"testing"

	"go4.org/xdgdir"
)

func TestGitManagerConfig_Open(t *testing.T) {
	gmc := GitManagerConfig{}

	fileName := filepath.Join(os.TempDir(), "gm.json")
	defaultFile := filepath.Join(xdgdir.Config.Path(), configPath)

	table := []struct {
		in  string
		out string
	}{
		{"", defaultFile},
		{fileName, filepath.Join(xdgdir.Config.Path(), fileName)},
	}

	for _, tt := range table {
		err := gmc.Open(tt.in)
		if err != nil {
			t.Fail()
			return
		}

		if gmc.file.Name() != tt.out {
			t.Errorf("got %q, want %q", gmc.file.Name(), tt.out)
		}

		if gmc.file == nil {
			t.Fail()
		}

	}

}
