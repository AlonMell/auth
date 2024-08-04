package main

import (
	"fmt"
	"log/slog"
	"os"
	"providerHub/internal/config"
	"providerHub/internal/storage/postgres"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))
	//log.Info(fmt.Sprintf("config: %+v", cfg))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg, log)
	if err != nil {
		log.Error("Error with start db postgres!" + err.Error())
	}

	log.Info(fmt.Sprintf("storage: %+v", storage))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
