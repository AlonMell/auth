package httpServer

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"providerHub/pkg/logger/sl"
	"time"
)

// Router is an interface that describes a router.
// It should be implemented by any router that is used in the server.
// The router should implement the [Convey] and [RegisterRoutes] methods.
// [Convey] is used to add middleware to the router.
// [RegisterRoutes] is used to register the routes.
type Router interface {
	http.Handler
	Convey()
	RegisterRoutes()
}

// Server is a wrapper around the http.Server.
type Server struct {
	*http.Server
	Logger *slog.Logger
}

// New creates a new server with the provided address and router.
// If the router is nil, a new [Mux] will be created.
func New(log *slog.Logger, address string, r Router) *Server {
	if r == nil {
		r = NewDefaultMux(log)
	}
	r.Convey()
	r.RegisterRoutes()

	return &Server{
		Server: &http.Server{
			Handler: r,
			Addr:    address,
		},
		Logger: log,
	}
}

// MustRun starts the server and blocks until it stops.
func (s *Server) MustRun() {
	s.Logger.Info("starting server", slog.String("address", s.Addr))
	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Logger.Error("server stop with error: ", sl.Err(err))
	}
}

// GracefulShutdown stops the server gracefully.
func (s *Server) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Logger.Info("shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Error("server forced to shutdown", sl.Err(err))
	}
}
