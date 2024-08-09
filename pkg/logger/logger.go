package logger

import (
	"log"
	"log/slog"
)

type TinyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type TinyHandler struct {
	opts slog.HandlerOptions
	slog.Handler
	l     *log.Logger
	attrs []slog.Attr
}
