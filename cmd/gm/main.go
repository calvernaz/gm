package main

import (
	"os"
	"path/filepath"
	
	"github.com/boltdb/bolt"
	"github.com/calvernaz/gm"
	"github.com/tidwall/tile38/controller/log"
)

const (
	configPath = "gm.json"
)

func main() {
	// load or create config file
	if err := loadOrCreateConfig(); err != nil {
		log.Fatalf("failed to create or load config: %v", err)
	}
}

func loadOrCreateConfig() interface{} {
	home := os.Getenv("HOME")
	
	// Open database.
	db, err := gm.Open(filepath.Join(home, configPath), 0666)
	if err != nil {
		return err
	}
	defer db.Close()
	
}
