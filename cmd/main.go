package main

import (
	"log/slog"
	"os"
	"os/signal"
	"providerHub/internal/config"
	"providerHub/internal/httpServer"
	"providerHub/internal/storage/postgres"
	"providerHub/pkg/logger"
	"providerHub/pkg/logger/sl"
	"syscall"
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
		log.Error("error with start db postgres!", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	server := httpServer.New(log, cfg.Address, nil)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go server.MustRun()

	<-stop
	server.GracefulShutdown()
	log.Info("server stopped")
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
