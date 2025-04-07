package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/AlonMell/grovelog"
	"github.com/AlonMell/grovelog/util"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = setupLogger(slog.LevelDebug, grovelog.Color)
	case envProd:
		log = setupLogger(slog.LevelWarn, grovelog.JSON)
	}

	return log
}

func setupLogger(level slog.Level, format grovelog.Format) *slog.Logger {
	opts := grovelog.NewOptions(level, "", format)
	return grovelog.NewLogger(os.Stdout, opts)
}

func WithLogOp(ctx context.Context, op string) context.Context {
	return util.UpdateLogCtx(ctx, "op", op)
}

func WithLogUserID(ctx context.Context, userID string) context.Context {
	return util.UpdateLogCtx(ctx, "userID", userID)
}

func WithLogRequestID(ctx context.Context, requestID string) context.Context {
	return util.UpdateLogCtx(ctx, "requestID", requestID)
}
