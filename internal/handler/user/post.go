package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlonMell/auth/internal/domain/dto"
	"github.com/AlonMell/auth/internal/handler/errors"
	resp "github.com/AlonMell/auth/internal/infra/lib/api/response"
	"github.com/AlonMell/auth/internal/infra/lib/decoder"
	"github.com/AlonMell/auth/internal/infra/lib/logger"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlonMell/auth/pkg/validator"
)

type Creater interface {
	Create(context.Context, dto.UserCreateReq) (id string, err error)
}

type PostResp struct {
	Id string `json:"id"`
}

func Post(
	log *slog.Logger, c Creater,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "handler.user.Post")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(ctx, log, w, r)

		var req dto.UserCreateReq
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.InfoContext(ctx, "request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		id, err := c.Create(ctx, req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.Status(r, http.StatusOK)
		resp.WriteJSON(w, r, PostResp{Id: id})
	}
}
