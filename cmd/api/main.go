package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Keyrozenxz/lessonhub/internal/config"
	"github.com/Keyrozenxz/lessonhub/internal/database"
	httptransport "github.com/Keyrozenxz/lessonhub/internal/transport/http"
)

func main() {
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

	router := httptransport.NewRouter()
	server := httptransport.NewServer(router, ":8080")

	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done() //Сервер ждет системный вызов, CTRL + C, Мейн блокируется, пока он не придет

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
}
