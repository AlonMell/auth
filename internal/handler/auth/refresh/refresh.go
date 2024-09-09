package refresh

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/domain/dto"
	"providerHub/internal/handler"
	"providerHub/internal/infra/config"
	resp "providerHub/internal/infra/lib/api/response"
)

type UserRefresher interface {
	RefreshToken(context.Context, dto.Refresh) (accessToken string, err error)
}

// New
// @Summary Refresh
// @Tags auth
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param input body Request true "refresh token"
// @Success 200 {object} Response
// @Router /refresh [post]
func New(
	log *slog.Logger, refresher UserRefresher, cfg config.JWT,
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

		req := Request{RefreshToken: cookie.Value}

		log.Info("request body decoded", slog.Any("request", req))

		refreshDTO := dto.Refresh{
			RefreshToken: req.RefreshToken,
			JWT:          cfg,
		}

		accessToken, err := refresher.RefreshToken(r.Context(), refreshDTO)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{AccessToken: accessToken})
	}
}
