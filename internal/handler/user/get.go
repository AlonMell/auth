package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlonMell/auth/internal/domain/dto"
	"github.com/AlonMell/auth/internal/handler/errors"
	resp "github.com/AlonMell/auth/internal/infra/lib/api/response"
	"github.com/AlonMell/auth/internal/infra/lib/logger"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlonMell/auth/internal/domain/model"
	"github.com/AlonMell/auth/pkg/validator"
)

type Getter interface {
	Get(context.Context, dto.UserGetReq) (*model.User, error)
}

type GetResp struct {
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	IsActive     bool   `json:"is_active"`
}

func Get(
	log *slog.Logger, g Getter,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "handler.user.Get")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(ctx, log, w, r)

		var req dto.UserGetReq
		req.Id = r.URL.Query().Get("id")

		log.InfoContext(ctx, "request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		ctx = logger.WithLogUserID(ctx, req.Id)

		u, err := g.Get(ctx, req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.Status(r, http.StatusOK)
		resp.WriteJSON(w, r, GetResp{u.Email, u.PasswordHash, u.IsActive})
	}
}
