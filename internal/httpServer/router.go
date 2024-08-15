package httpServer

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"providerHub/internal/httpServer/endpoint"
	mw "providerHub/internal/httpServer/middleware"
)

type Mux struct {
	*chi.Mux
	logger *slog.Logger
}

func NewDefaultMux(logger *slog.Logger) *Mux {
	return &Mux{
		Mux:    chi.NewRouter(),
		logger: logger,
	}
}

func (m *Mux) Convey() {
	m.Use(mw.Logger(m.logger))
	m.Use(mw.CORS)
	m.Use(middleware.RequestID)
	m.Use(middleware.Recoverer)
}

func (m *Mux) RegisterRoutes() {
	m.Get("/", endpoint.Home)
	m.Post("/register/", endpoint.Register)
	m.Get("/login/", endpoint.Login)
	m.Get("/api/isAdmin/", endpoint.IsAdmin)
}
