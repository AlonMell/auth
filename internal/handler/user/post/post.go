package post

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

type Creater interface {
	Create(context.Context, dto.UserCreate) (id string, err error)
}

// New
// @Summary Post User
// @Tags user
// @Security ApiKeyAuth
// @Description Create user at system
// @Accept json
// @Produce json
// @Param input body Request true "user info"
// @Success 200 {object} Response
// @Router /api/v1/user [post]
func New(
	log *slog.Logger, c Creater,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.post.New"

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

		createDTO := dto.UserCreate{
			Email:    req.Email,
			Password: req.Password,
			IsActive: req.IsActive,
		}

		id, err := c.Create(r.Context(), createDTO)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{Id: id})
	}
}
