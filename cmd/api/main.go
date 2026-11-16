package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	httptransport "github.com/Keyrozenxz/lessonhub/internal/transport/http"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := httptransport.NewRouter()
	server := httptransport.NewServer(router, ":8080")

	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done() //Сервер ждет системный вызов, CTRL + C, Мейн блокируется, пока он не придет

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
}
