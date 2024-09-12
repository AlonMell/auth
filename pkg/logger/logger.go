package logger

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	stdLog "log"
	"log/slog"
	"time"
)

type AttractiveHandlerOptions struct {
	SlogOpts      *slog.HandlerOptions
	TimeFormat    string
	UseAttractive bool
	UseFormat     bool
}

type AttractiveHandler struct {
	handler slog.Handler             //needed by main implement
	attrs   []slog.Attr              //needed by withAtrs
	opts    AttractiveHandlerOptions //needed by custom atrs and slogatrs
	l       *stdLog.Logger           //needed by handle out
}

type Message struct {
	Time  string
	Level string
	Msg   string
	Atrs  string
}

func NewAttractiveHandler(out io.Writer, opts AttractiveHandlerOptions) *AttractiveHandler {
	h := &AttractiveHandler{
		handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
		opts:    opts,
	}

	return h
}

func (h *AttractiveHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.opts.UseFormat {
		return h.handler.Handle(ctx, r)
	}

	timeStr := h.formatTime(r.Time)

	logMsg := r.Message
	formatLevel := r.Level.String() + ":"

	fields := h.collectFields(r)

	jsonOutput, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}
	output := string(jsonOutput)

	if !h.opts.UseAttractive {
		msg := Message{
			timeStr,
			formatLevel,
			logMsg,
			output,
		}
		return h.handle(ctx, msg)
	}

	type colorFn func(format string, a ...any) string
	levelColorMap := map[slog.Level]colorFn{
		slog.LevelDebug: color.MagentaString,
		slog.LevelInfo:  color.BlueString,
		slog.LevelWarn:  color.YellowString,
		slog.LevelError: color.RedString,
	}

	levelColorFunc := levelColorMap[r.Level]
	level := levelColorFunc(formatLevel)

	msg := Message{
		timeStr,
		level,
		color.CyanString(logMsg),
		color.WhiteString(output),
	}
	return h.handle(ctx, msg)
}

func (h *AttractiveHandler) handle(_ context.Context, msg Message) error {
	h.l.Println(msg.Time, msg.Level, msg.Msg, msg.Atrs)
	return nil
}

func (h *AttractiveHandler) formatTime(t time.Time) string {
	if h.opts.TimeFormat == "" {
		h.opts.TimeFormat = "[15:05:05.000]"
	}
	return t.Format(h.opts.TimeFormat)
}

func (h *AttractiveHandler) collectFields(r slog.Record) map[string]any {
	fields := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	return fields
}

func (h *AttractiveHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *AttractiveHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AttractiveHandler{
		handler: h.handler.WithAttrs(attrs),
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *AttractiveHandler) WithGroup(name string) slog.Handler {
	// TODO: implement
	return &AttractiveHandler{
		handler: h.handler.WithGroup(name),
		l:       h.l,
	}
}
