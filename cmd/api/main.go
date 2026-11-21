package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Keyrozenxz/lessonhub/internal/config"
	"github.com/Keyrozenxz/lessonhub/internal/database"
	httptransport "github.com/Keyrozenxz/lessonhub/internal/transport/http"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 5*time.Second)

	pool, err := database.NewPool(connectCtx, cfg)
	connectCancel()
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	logger.Info("connected to database", "host", cfg.DBHost, "port", cfg.DBPort, "name", cfg.DBName)

	addr := ":8080"

	router := httptransport.NewRouter()
	server := httptransport.NewServer(router, addr)

	go func() {
		logger.Info("starting http server", "addr", addr)
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done() //Сервер ждет системный вызов, CTRL + C, Мейн блокируется, пока он не придет

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	logger.Info("shutting down http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
	logger.Info("http server stopped")
}
