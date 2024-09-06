package get

import (
	"context"
	"log/slog"
	"net/http"
	"providerHub/internal/domain/dto"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/domain/model"
	"providerHub/internal/handler"
	resp "providerHub/internal/lib/api/response"
	"providerHub/pkg/validator"
)

type Getter interface {
	Get(context.Context, dto.UserGetDTO) (*model.User, error)
}

// New
// @Summary Get User
// @Tags user
// @Security ApiKeyAuth
// @Description Get User from system
// @Accept json
// @Produce json
// @Param input body Request true "user id"
// @Success 200 {object} Response
// @Router /api/v1/user/{id} [get]
func New(
	log *slog.Logger, g Getter,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.get.New"

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		req.Id = r.URL.Query().Get("id")

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		getDTO := dto.UserGetDTO{Id: req.Id}
		u, err := g.Get(r.Context(), getDTO)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{u.Email, u.PasswordHash, u.IsActive})
	}
}
