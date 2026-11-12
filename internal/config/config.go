package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" || portStr == "" || user == "" || password == "" || name == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("DB_PORT must be between 1 and 65535")
	}

	return &Config{
		DBHost:     host,
		DBPort:     port,
		DBUser:     user,
		DBName:     name,
		DBPassword: password,
	}, nil
}
