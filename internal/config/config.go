package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL           string
	Port                  string
	TrueFoundryAPIKey     string
	TrueFoundryBaseURL    string
	TrueFoundryCheapModel string
	TrueFoundryExpModel   string
	AerospikeHost         string
	AerospikePort         string
	AerospikeNamespace    string
}

func Load() (*Config, error) {
	// Load .env file if it exists (development). Ignore error if not found.
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		Port:                  os.Getenv("PORT"),
		TrueFoundryAPIKey:     os.Getenv("TRUEFOUNDRY_API_KEY"),
		TrueFoundryBaseURL:    os.Getenv("TRUEFOUNDRY_BASE_URL"),
		TrueFoundryCheapModel: os.Getenv("TRUEFOUNDRY_CHEAP_MODEL"),
		TrueFoundryExpModel:   os.Getenv("TRUEFOUNDRY_EXPENSIVE_MODEL"),
		AerospikeHost:         os.Getenv("AEROSPIKE_HOST"),
		AerospikePort:         os.Getenv("AEROSPIKE_PORT"),
		AerospikeNamespace:    os.Getenv("AEROSPIKE_NAMESPACE"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL env var is required")
	}
	return cfg, nil
}
