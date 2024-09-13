package user

import (
	"context"
	"github.com/AlonMell/ProviderHub/internal/delivery/http/catcher"
	"github.com/AlonMell/ProviderHub/internal/domain/dto"
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/decoder"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlonMell/ProviderHub/pkg/validator"
)

type Deleter interface {
	Delete(context.Context, dto.UserDeleteReq) error
}

type DeleteResp struct {
	resp.Response
}

// Delete
// @Summary Delete User
// @Tags user
// @Security ApiKeyAuth
// @Description Delete User from system
// @Accept json
// @Produce json
// @Param input body Request true "user id"
// @Success 200 {object} Response
// @Router /api/v1/user [delete]
func Delete(
	log *slog.Logger, d Deleter,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "http.user.Delete")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := catcher.NewCatcher(ctx, log, w, r)

		var req dto.UserDeleteReq
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.InfoContext(ctx, "request body decoded", slog.Any("request", req))
		ctx = logger.WithLogUserID(ctx, req.Id)

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		if err := d.Delete(ctx, req); err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.Status(r, http.StatusOK)
		resp.WriteJSON(w, r, DeleteResp{Response: resp.Ok()})
	}
}
