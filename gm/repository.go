package gm

import (
	"time"
)

// Repository format entry
type Repository struct {
	// Repository Name
	Name string `json:"name"`
	// Tracking repository
	Enabled bool `json:"enabled"`
	// Last repository update
	LastUpdate time.Time `json:"last_update"`
	// Git Repository path
	Path string
}
