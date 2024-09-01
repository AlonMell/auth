package logger

import (
	"log/slog"
	"os"

	"providerHub/pkg/logger"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func SetupLogger(env string) *slog.Logger {
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
