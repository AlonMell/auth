package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/router/handler/auth"
	"providerHub/internal/router/handler/auth/dto"
	mw "providerHub/internal/router/middleware"
)

// Router is an interface that describes a router.
// It should be implemented by any router that is used in the server.
// The router should implement the [Convey] and [RegisterRoutes] methods.
// [Convey] is used to add middleware to the router.
// [HandleAuth] is used to add authentication to the router.
type Router interface {
	http.Handler
	Convey()
	HandleAuth()
}

// Auth is an interface that describes the authentication service.
type Auth interface {
	RegisterUser(dto.RegisterRequest) (userId string, err error)
	Token(dto.LoginRequest) (token string, err error)
}

type Mux struct {
	*chi.Mux
	Auth   Auth
	logger *slog.Logger
}

func New(logger *slog.Logger, auth Auth) *Mux {
	return &Mux{
		Mux:    chi.NewRouter(),
		logger: logger,
		Auth:   auth,
	}
}

func (m *Mux) Prepare() {
	m.Convey()
	m.HandleAuth()
}

func (m *Mux) Convey() {
	m.Use(mw.Logger(m.logger))
	m.Use(mw.CORS)
	m.Use(middleware.RequestID)
	m.Use(middleware.Recoverer)
}

func (m *Mux) HandleAuth() {
	m.Get("/", auth.Main)
	m.Post("/register", auth.Register(m.logger, m.Auth))
	m.Get("/login", auth.Login(m.logger, m.Auth))
}
