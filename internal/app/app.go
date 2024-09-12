package app

import (
	httpApp "github.com/AlonMell/ProviderHub/internal/app/http"
	"github.com/AlonMell/ProviderHub/internal/app/postgres"
	"github.com/AlonMell/ProviderHub/internal/infra/config"
	"github.com/AlonMell/ProviderHub/internal/infra/repo"
	"github.com/AlonMell/ProviderHub/internal/router"
	"github.com/AlonMell/ProviderHub/internal/service/auth"
	"github.com/AlonMell/ProviderHub/internal/service/user"
	"github.com/AlonMell/ProviderHub/pkg/logger/sl"
	sq "github.com/Masterminds/squirrel"
	"log/slog"
	"os"
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
