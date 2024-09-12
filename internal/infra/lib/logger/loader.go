package logger

import (
	"log/slog"
	"os"

	"github.com/AlonMell/ProviderHub/pkg/logger"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

//TODO: Сделать интерфейс
/*type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}*/

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = setupAttractiveLogger(slog.LevelDebug, true, true)
	case envProd:
		log = setupAttractiveLogger(slog.LevelWarn, false, false)
	}

	return log
}

func setupAttractiveLogger(
	level slog.Level, useAttractive bool, useFormat bool,
) *slog.Logger {
	opts := logger.AttractiveHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
		UseAttractive: useAttractive,
		UseFormat:     useFormat,
	}

	handler := logger.NewAttractiveHandler(os.Stdout, opts)

	developHandler := NewDevelopHandler(handler)

	return slog.New(developHandler)
}
