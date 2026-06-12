package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Theme string `json:"theme"`
}

func LoadConfig() Config {
	cfg := Config{Theme: "dark"}
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}
	configPath := filepath.Join(home, ".pi", "agent", "pi-mc.json")
	b, err := os.ReadFile(configPath)
	if err == nil {
		_ = json.Unmarshal(b, &cfg)
	}
	return cfg
}
