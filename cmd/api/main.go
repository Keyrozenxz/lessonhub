package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Keyrozenxz/lessonhub/internal/config"
	"github.com/Keyrozenxz/lessonhub/internal/database"
	httptransport "github.com/Keyrozenxz/lessonhub/internal/transport/http"
)

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 5*time.Second)

	pool, err := database.NewPool(connectCtx, cfg)
	connectCancel()
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}
	defer pool.Close()
	logger.Info("connected to database", "host", cfg.DBHost, "port", cfg.DBPort, "name", cfg.DBName)

	addr := ":8080"

	router := httptransport.NewRouter()
	server := httptransport.NewServer(router, addr)

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting http server", "addr", addr)
		if err := server.Start(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case v := <-serverErr:
		return fmt.Errorf("run http server: %w", v)
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	logger.Info("shutting down http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	logger.Info("http server stopped")

	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("application stopped with error", "err", err)
		os.Exit(1)
	}

}
