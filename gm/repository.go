package gm

import (
	"encoding/json"
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

// UnmarshalJSON is needed because Parsed has unexported fields.
func (r *Repository) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r)
}

// MarshalJSON is needed because Parsed has unexported fields.
func (r *Repository) MarshalJSON() ([]byte, error) {
	return json.Marshal(&r)
}
