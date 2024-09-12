package delete

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

type Deleter interface {
	Delete(context.Context, dto.UserDelete) error
}

// New
// @Summary Delete User
// @Tags user
// @Security ApiKeyAuth
// @Description Delete User from system
// @Accept json
// @Produce json
// @Param input body Request true "user id"
// @Success 200 {object} Response
// @Router /api/v1/user [delete]
func New(
	log *slog.Logger, d Deleter,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.delete.New"
		ctx := logger.WithLogOp(r.Context(), op)
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(op, log, w, r)

		var req Request
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		ctx = logger.WithLogUserID(ctx, req.Id)

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		deleteDTO := dto.UserDelete{Id: req.Id}

		if err := d.Delete(ctx, deleteDTO); err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{Response: resp.Ok()})
	}
}
