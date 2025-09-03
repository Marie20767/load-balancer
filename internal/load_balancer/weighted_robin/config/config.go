package config

import (
	"encoding/json"
	"os"
)

type Server struct {
	URL    string
	Weight int
}

type ServerConfig struct {
	Servers []Server `json:"servers"`
}

func LoadConfig() ([]Server, error) {
	f, err := os.Open("robin.config.json")

	if err != nil {
		return nil, err
	}

	defer f.Close()

	var cfg ServerConfig

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg.Servers, nil
}

