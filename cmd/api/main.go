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
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	connectCtx, connectCancel := context.WithTimeout(signalCtx, 5*time.Second)

	pool, err := database.NewPool(connectCtx, cfg.DB)
	connectCancel()
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}
	defer pool.Close()
	logger.Info("connected to database", "host", cfg.DB.Host, "port", cfg.DB.Port, "name", cfg.DB.Name)

	router := httptransport.NewRouter()
	server := httptransport.NewServer(router, cfg.HTTP.Addr)

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting http server", "addr", cfg.HTTP.Addr)
		if err := server.Start(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("run http server: %w", err)
	case <-signalCtx.Done():
		logger.Info("shutdown signal received")
		stop()
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
