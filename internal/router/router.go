package router

import (
	"github.com/AlonMell/ProviderHub/internal/infra/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"

	_ "github.com/AlonMell/ProviderHub/api"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/AlonMell/ProviderHub/internal/handler/auth/refresh"

	"github.com/AlonMell/ProviderHub/internal/handler/user/get"
	mw "github.com/AlonMell/ProviderHub/internal/router/middleware"

	"github.com/AlonMell/ProviderHub/internal/handler/auth"
	"github.com/AlonMell/ProviderHub/internal/handler/auth/login"
	"github.com/AlonMell/ProviderHub/internal/handler/auth/register"

	"github.com/AlonMell/ProviderHub/internal/handler/user"
	del "github.com/AlonMell/ProviderHub/internal/handler/user/delete"
	"github.com/AlonMell/ProviderHub/internal/handler/user/post"
	"github.com/AlonMell/ProviderHub/internal/handler/user/update"
)

type Router interface {
	http.Handler
	Convey()
	HandleAuth(config.JWT)
	HandleUsers(config.JWT)
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

func (m *Mux) Prepare(cfg config.JWT) {
	m.Convey()
	m.HandleUsers(cfg)
	m.HandleAuth(cfg)
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

func (m *Mux) HandleUsers(cfg config.JWT) {
	m.With(mw.Auth(cfg)).Get("/api/v1/users/{id}", get.New(m.logger, m.UserProvider))
	m.With(mw.Auth(cfg)).Post("/api/v1/users", post.New(m.logger, m.UserProvider))
	m.With(mw.Auth(cfg)).Put("/api/v1/users", update.New(m.logger, m.UserProvider))
	m.With(mw.Auth(cfg)).Delete("/api/v1/users", del.New(m.logger, m.UserProvider))
}

func (m *Mux) HandleAuth(cfg config.JWT) {
	m.Post("/register", register.New(m.logger, m.Auth))
	m.Post("/login", login.New(m.logger, cfg, m.Auth))
	m.Post("/refresh", refresh.New(m.logger, m.Auth, cfg))
}
