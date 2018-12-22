package internal

import (
	"testing"
)

func TestVcsDownload(t *testing.T) {
	vcsGit.download("/Users/calvernaz/sandbox/go-repository/public/src/gobot.io/x/gobot")
}
