package http

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(handler http.Handler, addr string) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &Server{server: srv}
}
func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	return nil
}
