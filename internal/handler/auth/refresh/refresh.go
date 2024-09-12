package refresh

import (
	"context"
	"github.com/AlonMell/ProviderHub/internal/domain/dto"
	"github.com/AlonMell/ProviderHub/internal/handler/errors"
	"github.com/AlonMell/ProviderHub/internal/infra/config"
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
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
		ctx := logger.WithLogOp(r.Context(), op)
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(op, log, w, r)

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

		accessToken, err := refresher.RefreshToken(ctx, refreshDTO)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{AccessToken: accessToken})
	}
}
