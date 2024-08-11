package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// CORS принимает параметром Handler и возвращает тоже Handler.
func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// замыкание: используем ServeHTTP следующего хендлера
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// log *slog.Logger
func Logger(next http.Handler) http.Handler {
	log := log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")

	fn := func(w http.ResponseWriter, r *http.Request) {
		//r.Context()
		entry := log.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			entry.Info("request completed",
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("duration", time.Since(t1).String()),
			)
		}()

		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
