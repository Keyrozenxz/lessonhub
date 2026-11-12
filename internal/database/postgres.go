package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/Keyrozenxz/lessonhub/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   net.JoinHostPort(cfg.DBHost, strconv.Itoa(cfg.DBPort)),
		Path:   "/" + cfg.DBName,
	}

	connString := u.String()

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
