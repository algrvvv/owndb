package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DebugMode      bool          `yaml:"debug"`
	SnapshotFile   string        `yaml:"snapshot_file"`
	WalFile        string        `yaml:"wal_file"`
	LogFile        string        `yaml:"log_file"`
	DumperInterval time.Duration `yaml:"dumper_interval"`
	Port           int           `yaml:"port"`
}

func MustLoad() *Config {
	path := configPath()

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
