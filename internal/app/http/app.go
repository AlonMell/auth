package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlonMell/ProviderHub/internal/infra/config"
	"log/slog"
	"net/http"
	"time"

	"github.com/AlonMell/ProviderHub/internal/router"
	"github.com/AlonMell/ProviderHub/pkg/logger/sl"
)

type Server struct {
	*http.Server
	Logger *slog.Logger
}

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

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	const op = "httpServer.Run"

	s.Logger.Info("starting server", slog.String("address", s.Addr))

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Logger.Info("shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Error("server forced to shutdown", sl.Err(err))
	}
}
