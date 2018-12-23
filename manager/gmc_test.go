package manager

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGitManagerConfig_Open(t *testing.T) {
	gmc := GitManagerConfig{}
	
	f, _ := ioutil.TempFile(os.TempDir(), "gm.json")
	
	table := []struct {
		in string
		out string
	} {
		{"", "gm/gm.json"},
		{f.Name(), f.Name()},
	}
	
	for _, tt := range table {
		err := gmc.Open(tt.in)
		if err != nil {
			t.Fail()
			return
		}
		
		if gmc.path != tt.out {
			t.Errorf("got %q, want %q", gmc.path, tt.out)
		}
		
		if gmc.file == nil {
			t.Fail()
		}
		
	}
	
}
