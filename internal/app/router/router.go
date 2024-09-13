package router

import (
	mw "github.com/AlonMell/ProviderHub/internal/app/router/middleware"
	"github.com/AlonMell/ProviderHub/internal/delivery/http/auth"
	"github.com/AlonMell/ProviderHub/internal/delivery/http/user"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/AlonMell/ProviderHub/api"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router interface {
	http.Handler
	Convey()
	HandleAuth(refreshTTl time.Duration)
	HandleUsers(cfg jwt.Config)
}

type Mux struct {
	*chi.Mux
	Auth         auth.Auth
	UserProvider user.Provider
	logger       *slog.Logger
}

func New(
	logger *slog.Logger, auth auth.Auth, provider user.Provider,
) *Mux {
	return &Mux{
		Mux:          chi.NewRouter(),
		logger:       logger,
		Auth:         auth,
		UserProvider: provider,
	}
}

func (m *Mux) Prepare(cfg jwt.Config) {
	m.Convey()
	m.HandleUsers(cfg)
	m.HandleAuth(cfg.RefreshTTL)
	m.Swagger()
}

func (m *Mux) Convey() {
	m.Use(middleware.RequestID)
	m.Use(mw.CORS)
	m.Use(middleware.URLFormat)
	m.Use(mw.Logger(m.logger))
	m.Use(middleware.Recoverer)
}

func (m *Mux) Swagger() {
	m.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}

func (m *Mux) HandleUsers(cfg jwt.Config) {
	m.With(mw.Auth(cfg)).Get("/api/v1/users/{id}", user.Get(m.logger, m.UserProvider))
	m.With(mw.Auth(cfg)).Post("/api/v1/users", user.Post(m.logger, m.UserProvider))
	m.With(mw.Auth(cfg)).Put("/api/v1/users", user.Update(m.logger, m.UserProvider))
	m.With(mw.Auth(cfg)).Delete("/api/v1/users", user.Delete(m.logger, m.UserProvider))
}

func (m *Mux) HandleAuth(refreshTTl time.Duration) {
	m.Post("/register", auth.Register(m.logger, m.Auth))
	m.Post("/login", auth.Login(m.logger, refreshTTl, m.Auth))
	m.Post("/refresh", auth.Refresh(m.logger, m.Auth))
}
