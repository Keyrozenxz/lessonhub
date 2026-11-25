package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type HTTPConfig struct {
	Addr string
}

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	httpAddr := os.Getenv("HTTP_ADDR")

	if dbHost == "" || dbPortStr == "" || dbUser == "" || dbPassword == "" || dbName == "" || httpAddr == "" {
		return nil, fmt.Errorf("configuration is incomplete")
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	if dbPort < 1 || dbPort > 65535 {
		return nil, fmt.Errorf("DB_PORT must be between 1 and 65535")
	}

	return &Config{
		DB: DBConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Name:     dbName,
			Password: dbPassword,
		},
		HTTP: HTTPConfig{Addr: httpAddr},
	}, nil
}
