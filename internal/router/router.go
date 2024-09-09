package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/infra/config"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "providerHub/api"

	"providerHub/internal/handler/auth/refresh"

	"providerHub/internal/handler/user/get"
	mw "providerHub/internal/router/middleware"

	"providerHub/internal/handler/auth"
	"providerHub/internal/handler/auth/login"
	"providerHub/internal/handler/auth/register"

	"providerHub/internal/handler/user"
	del "providerHub/internal/handler/user/delete"
	"providerHub/internal/handler/user/post"
	"providerHub/internal/handler/user/update"
)

type Router interface {
	http.Handler
	Convey()
	HandleAuth(config.JWT)
	HandleUsers(config.JWT)
}

type Mux struct {
	*chi.Mux
	Auth         auth.Interface
	UserProvider user.Interface
	logger       *slog.Logger
}

func New(
	logger *slog.Logger, auth auth.Interface, provider user.Interface,
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
