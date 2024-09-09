package delete

import (
	"context"
	"log/slog"
	"net/http"
	"providerHub/internal/domain/dto"
	resp "providerHub/internal/infra/lib/api/response"
	"providerHub/internal/infra/lib/decoder"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/handler"
	"providerHub/pkg/validator"
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

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		deleteDTO := dto.UserDelete{Id: req.Id}

		if err := d.Delete(r.Context(), deleteDTO); err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{Response: resp.Ok()})
	}
}
