package logger

import (
	"context"
	"errors"
	"log/slog"
)

type ctxKeyLog int

const LogKey = ctxKeyLog(0)

type logCtx map[string]string

type DevelopHandler struct {
	handler slog.Handler
}

func NewDevelopHandler(h slog.Handler) *DevelopHandler {
	return &DevelopHandler{handler: h}
}

func (l *DevelopHandler) Enabled(ctx context.Context, lev slog.Level) bool {
	return l.handler.Enabled(ctx, lev)
}

//slog string?

func (l *DevelopHandler) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := getLogCtx(ctx); ok {
		for key, value := range c {
			rec.Add(key, value)
		}
	}
	return l.handler.Handle(ctx, rec)
}

func (l *DevelopHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DevelopHandler{handler: l.handler.WithAttrs(attrs)}
}

func (l *DevelopHandler) WithGroup(name string) slog.Handler {
	return &DevelopHandler{handler: l.handler.WithGroup(name)}
}

func updateLogCtx(ctx context.Context, newCtx logCtx) context.Context {
	if existingCtx, ok := getLogCtx(ctx); ok {
		for k, v := range newCtx {
			existingCtx[k] = v
		}
		return context.WithValue(ctx, LogKey, existingCtx)
	}
	return context.WithValue(ctx, LogKey, newCtx)
}

func getLogCtx(ctx context.Context) (logCtx, bool) {
	c, ok := ctx.Value(LogKey).(logCtx)
	return c, ok
}

func WithLogUserID(ctx context.Context, userID string) context.Context {
	return updateLogCtx(ctx, logCtx{"userID": userID})
}

func WithLogOp(ctx context.Context, op string) context.Context {
	if existingCtx, ok := getLogCtx(ctx); ok {
		existingCtx["op"] += ": " + op
		return context.WithValue(ctx, LogKey, existingCtx)
	}
	return context.WithValue(ctx, LogKey, logCtx{"op": op})
}

func WithLogRequestID(ctx context.Context, requestID string) context.Context {
	return updateLogCtx(ctx, logCtx{"requestID": requestID})
}

type errorWithLogCtx struct {
	err error
	ctx logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.err.Error()
}

func (e *errorWithLogCtx) Wrap(ctx context.Context, err error) error {
	c, _ := getLogCtx(ctx)
	return &errorWithLogCtx{err: err, ctx: c}
}

func (e *errorWithLogCtx) ErrorCtx(ctx context.Context, err error) context.Context {
	var errCtx *errorWithLogCtx
	if errors.As(err, &errCtx) {
		return updateLogCtx(ctx, errCtx.ctx)
	}
	return ctx
}
