package refresh

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/config"
	"providerHub/internal/handler"
	"providerHub/internal/handler/auth"
	resp "providerHub/internal/lib/api/response"
)

type UserRefresher interface {
	RefreshToken(auth.RefreshRequest, config.JWT) (accessToken string, err error)
}

func New(
	log *slog.Logger,
	refresher UserRefresher,
	cfg config.JWT,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.refresh.New"

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		cookie, err := r.Cookie("refresh")
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		req := auth.RefreshRequest{RefreshToken: cookie.Value}

		log.Info("request body decoded", slog.Any("request", req))

		accessToken, err := refresher.RefreshToken(req, cfg)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.WriteJSON(w, r, auth.RefreshResponse{AccessToken: accessToken})
	}
}
