package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	LogLevel    string

	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_SSLMODE  string
}

func Load() (*Config, error) {
	cfg := Defaults

	godotenv.Load() // optional file, ignore if missing

	if v := os.Getenv("APP_PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.Environment = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}

	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		cfg.POSTGRES_HOST = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		cfg.POSTGRES_PORT = v
	}
	if v := os.Getenv("POSTGRES_USER"); v != "" {
		cfg.POSTGRES_USER = v
	}
	if v := os.Getenv("POSTGRES_PASSWORD"); v != "" {
		cfg.POSTGRES_PASSWORD = v
	}
	if v := os.Getenv("POSTGRES_DB"); v != "" {
		cfg.POSTGRES_DB = v
	}
	if v := os.Getenv("POSTGRES_SSLMODE"); v != "" {
		cfg.POSTGRES_SSLMODE = v
	}

	return &cfg, nil
}
