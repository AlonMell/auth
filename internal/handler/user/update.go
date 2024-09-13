package user

import (
	"context"
	"github.com/AlonMell/ProviderHub/internal/domain/dto"
	"github.com/AlonMell/ProviderHub/internal/handler/errors"
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/decoder"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlonMell/ProviderHub/pkg/validator"
)

type Updater interface {
	Update(context.Context, dto.UserUpdateReq) error
}

type UpdateResp struct {
	resp.Response
}

// Update
// @Summary Update User
// @Tags user
// @Security ApiKeyAuth
// @Description Update all info of a user at system
// @Accept json
// @Produce json
// @Param input body Request true "all user info"
// @Success 200 {object} Response
// @Router /api/v1/user [put]
func Update(
	log *slog.Logger, u Updater,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "handler.user.Update")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(ctx, log, w, r)

		var req dto.UserUpdateReq
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		if err := u.Update(ctx, req); err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, UpdateResp{resp.Ok()})
	}
}
