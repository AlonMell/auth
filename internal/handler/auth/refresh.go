package auth

import (
	"context"
	"github.com/AlonMell/ProviderHub/internal/domain/dto"
	"github.com/AlonMell/ProviderHub/internal/handler/errors"
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type UserRefresher interface {
	RefreshToken(context.Context, dto.RefreshReq) (accessToken string, err error)
}

type RefreshResp struct {
	AccessToken string `json:"access_token"`
}

// Refresh
// @Summary Refresh
// @Tags auth
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param input body Request true "refresh token"
// @Success 200 {object} Response
// @Router /refresh [post]
func Refresh(
	log *slog.Logger, refresher UserRefresher,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "handler.auth.Refresh")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(ctx, log, w, r)

		cookie, err := r.Cookie("refresh")
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		req := dto.RefreshReq{RefreshToken: cookie.Value}

		log.Info("request body decoded", slog.Any("request", req))

		accessToken, err := refresher.RefreshToken(ctx, req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.Status(r, http.StatusOK)
		resp.WriteJSON(w, r, RefreshResp{AccessToken: accessToken})
	}
}
