package get

import (
	"context"
	"github.com/AlonMell/ProviderHub/internal/domain/dto"
	"github.com/AlonMell/ProviderHub/internal/handler/errors"
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlonMell/ProviderHub/internal/domain/model"
	"github.com/AlonMell/ProviderHub/pkg/validator"
)

type Getter interface {
	Get(context.Context, dto.UserGet) (*model.User, error)
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
		ctx := logger.WithLogOp(r.Context(), op)
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(op, log, w, r)

		var req Request
		req.Id = r.URL.Query().Get("id")

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		getDTO := dto.UserGet{Id: req.Id}

		ctx = logger.WithLogUserID(ctx, getDTO.Id)

		u, err := g.Get(ctx, getDTO)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{u.Email, u.PasswordHash, u.IsActive})
	}
}
