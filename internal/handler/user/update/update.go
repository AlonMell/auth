package update

import (
	"context"
	"log/slog"
	"net/http"
	"providerHub/internal/domain/dto"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/handler"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/validator"
)

type Updater interface {
	Update(context.Context, dto.UserUpdateDTO) error
}

// New
// @Summary Update User
// @Tags user
// @Security ApiKeyAuth
// @Description Update all info of a user at system
// @Accept json
// @Produce json
// @Param input body Request true "all user info"
// @Success 200 {object} Response
// @Router /api/v1/user [put]
func New(
	log *slog.Logger, u Updater,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.update.New"

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

		updateDTO := dto.UserUpdateDTO{
			Email:    req.Email,
			Password: req.Password,
			IsActive: req.IsActive,
		}

		if err := u.Update(r.Context(), updateDTO); err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{resp.Ok()})
	}
}
