package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port         int    `env:"PORT" envDefault:"8080"`
	Env          string `env:"ENV" envDefault:"development"`
	DBConnString string `env:"DB_CONN_STRING,required"`
	JWTSecret    string `env:"JWT_SECRET,required"`
}

func Load() (*Config, error) {
	// Load .env file if it exists (useful for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
