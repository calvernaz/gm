package gm

import (
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	repo := Repository{
		Name:       "test",
		Enabled:    true,
		LastUpdate: time.Now(),
		Path:       "/Users/calvernaz/sandbox/go-repository/public/src/github.com/gorilla/sessions",
	}
	
	repo.Update()
}
