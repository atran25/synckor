package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log/slog"
)

type Config struct {
	Port                int    `env:"PORT"`
	RegistrationEnabled bool   `env:"REGISTRATION_ENABLED"`
	AdminUsername       string `env:"ADMIN_USERNAME"`
	AdminPassword       string `env:"ADMIN_PASSWORD"`
}

func LoadConfig() (Config, error) {
	// Load configuration from file
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Failed to load configuration file. Does it exist?", "Error", err)
	}
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		slog.Error("Failed to parse environment variable into Config struct.", "Error", err)
		return Config{}, err
	}
	return cfg, nil
}
