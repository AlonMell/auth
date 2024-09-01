package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"providerHub/internal/config"
	"providerHub/internal/router"
	"providerHub/pkg/logger/sl"
)

// Server is a struct that represents an HTTP server.
// It embeds the [http.Server] struct and adds a logger.
type Server struct {
	*http.Server
	Logger *slog.Logger
}

// New creates a new server with the provided address and router.
// If the router is nil, a new [Mux] will be created.
func New(log *slog.Logger, cfg config.HTTPServer, r router.Router) *Server {
	return &Server{
		Server: &http.Server{
			Handler:      r,
			Addr:         cfg.Address,
			WriteTimeout: cfg.Timeout,
			ReadTimeout:  cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		Logger: log,
	}
}

// MustRun starts the server and panic if an error occurs.
// It is a wrapper around the [Run] method.
func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

// Run starts the server.
// If the server stops, it returns an error.
func (s *Server) Run() error {
	const op = "httpServer.Run"

	s.Logger.Info("starting server", slog.String("address", s.Addr))

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops the server gracefully.
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Logger.Info("shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Error("server forced to shutdown", sl.Err(err))
	}
}
