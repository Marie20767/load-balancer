package config

import (
	"encoding/json"
	"os"
)

type Server struct {
	URL      string
	Position float32
}

type ServerConfig struct {
	Servers []Server `json:"servers"`
}

const HashRingValue = 1

func LoadConfig() ([]Server, error) {
	f, err := os.Open("hashing.config.json")

	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg ServerConfig

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	servers := cfg.Servers

	hashRingPortion := HashRingValue / float32(len(cfg.Servers))

	for i := range servers {
		servers[i].Position = float32((i + 1)) * hashRingPortion
	}

	return servers, nil
}
