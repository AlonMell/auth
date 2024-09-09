package app

import (
	sq "github.com/Masterminds/squirrel"
	"log/slog"
	"os"
	httpApp "providerHub/internal/app/http"
	"providerHub/internal/app/postgres"
	"providerHub/internal/infra/config"
	"providerHub/internal/infra/repo"
	"providerHub/internal/router"
	"providerHub/internal/service/auth"
	"providerHub/internal/service/user"
	"providerHub/pkg/logger/sl"
)

var (
	postgresPlaceholder = sq.Dollar
)

type App struct {
	Server *httpApp.Server
}

func New(log *slog.Logger, cfg *config.Config) *App {
	db, err := postgres.New(cfg, log)
	if err != nil {
		log.Error("error with start db postgres!", sl.Err(err))
		os.Exit(1)
	}

	userRepo := repo.NewUserRepo(db, postgresPlaceholder)

	authService := auth.New(log, userRepo)
	userService := user.New(log, userRepo)

	mux := router.New(log, authService, userService)
	mux.Prepare(cfg.JWT)

	server := httpApp.New(log, cfg.HTTPServer, mux)

	return &App{Server: server}
}
