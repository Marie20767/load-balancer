package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func ParseEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, errors.New("port environment variable not set")
	}

	return &Config{
		Port: port,
	}, nil
}

type Server struct {
	URL    string
	Weight int
}

type ServerConfig struct {
	Servers []Server `json:"servers"`
}

func LoadConfig() ([]Server, error) {
	f, err := os.Open("weighted_robin.config.json")

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

