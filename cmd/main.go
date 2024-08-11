package main

import (
	"log/slog"
	"os"
	"providerHub/internal/config"
	"providerHub/internal/httpServer"
	"providerHub/internal/storage/postgres"
	"providerHub/pkg/logger"
	"providerHub/pkg/logger/sl"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.Any("cfg", cfg))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg, log)
	if err != nil {
		log.Error("Error with start db postgres!", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	err = httpServer.Run()
	if err != nil {
		log.Error("Error with start http server:", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = setupDevelopLogger()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupDevelopLogger() *slog.Logger {
	opts := logger.DevelopHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewDevelopHandler(os.Stdout)

	return slog.New(handler)
}
