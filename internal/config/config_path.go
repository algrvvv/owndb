package config

import (
	"flag"
	"os"
	"path/filepath"
)

var (
	path                 string
	envConfigPathVarName = "OWN_CONFIG"
)

func Path() string {
	flag.StringVar(&path, "config", "", "config path")
	flag.Parse()

	if path == "" {
		path = os.Getenv(envConfigPathVarName)
	}

	if path == "" {
		path = defaultPath()
	}

	return path
}

func defaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".owndb/config.yml")
}
