package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	var missing []string

	if dbHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if dbPortStr == "" {
		missing = append(missing, "DB_PORT")
	}
	if dbUser == "" {
		missing = append(missing, "DB_USER")
	}
	if dbPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if dbName == "" {
		missing = append(missing, "DB_NAME")
	}
	if httpAddr == "" {
		missing = append(missing, "HTTP_ADDR")
	}

	if len(missing) > 0 {
		missingStr := strings.Join(missing, ", ")
		return nil, fmt.Errorf("missing required environment variables: %s", missingStr)
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
