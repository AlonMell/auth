package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/domain/model"
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

// Router is an interface that describes a router.
// It should be implemented by any router that is used in the server.
// The router should implement the [Convey] and [RegisterRoutes] methods.
// [Convey] is used to add middleware to the router.
// [HandleAuth] is used to add authentication to the router.
type Router interface {
	http.Handler
	Convey()
	HandleAuth(tokenTTL time.Duration)
	HandleUsers()
}

// Auth is an interface that describes the authentication service.
type Auth interface {
	RegisterUser(auth.RegisterRequest) (userId string, err error)
	Token(auth.LoginRequest, time.Duration) (token string, err error)
}

// UserProvider is an interface that describes the user service.
type UserProvider interface {
	Get(user.GetUserRequest) (*model.User, error)
	Create(user.CreateUserRequest) (string, error)
	Delete(user.DeleteUserRequest) error
	Update(user.UpdateUserRequest) error
}

type Mux struct {
	*chi.Mux
	Auth         Auth
	UserProvider UserProvider
	logger       *slog.Logger
}

func New(logger *slog.Logger, auth Auth, provider UserProvider) *Mux {
	return &Mux{
		Mux:          chi.NewRouter(),
		logger:       logger,
		Auth:         auth,
		UserProvider: provider,
	}
}

func (m *Mux) Prepare(tokenTTL time.Duration) {
	m.Convey()
	m.HandleUsers()
	m.HandleAuth(tokenTTL)
}

func (m *Mux) Convey() {
	m.Use(mw.Logger(m.logger))
	m.Use(mw.CORS)
	m.Use(middleware.RequestID)
	m.Use(middleware.Recoverer)
	m.Use(middleware.URLFormat)
}

func (m *Mux) HandleUsers() {
	m.Get("/api/v1/users", get.New(m.logger, m.UserProvider))
	m.Post("/api/v1/users", post.New(m.logger, m.UserProvider))
	m.Put("/api/v1/users", update.New(m.logger, m.UserProvider))
	m.Delete("/api/v1/users", del.New(m.logger, m.UserProvider))
}

func (m *Mux) HandleAuth(tokenTTL time.Duration) {
	m.Post("/register", register.New(m.logger, m.Auth))
	m.Post("/login", login.New(m.logger, tokenTTL, m.Auth))
}
