package app

import (
	"log/slog"
	"os"
	httpApp "providerHub/internal/app/http"
	"providerHub/internal/config"
	"providerHub/internal/router"
	"providerHub/internal/service/auth"
	"providerHub/internal/storage/postgres"
	"providerHub/pkg/logger/sl"
)

type App struct {
	Server *httpApp.Server
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := postgres.New(cfg, log)
	if err != nil {
		log.Error("error with start db postgres!", sl.Err(err))
		os.Exit(1)
	}

	authService := auth.New(log, storage, storage)

	mux := router.New(log, authService)
	mux.Prepare()

	server := httpApp.New(log, cfg.Address, mux)

	return &App{Server: server}
}
