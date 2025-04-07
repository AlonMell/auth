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

type Deleter interface {
	Delete(context.Context, dto.UserDeleteReq) error
}

type DeleteResp struct {
	resp.Response
}

func Delete(
	log *slog.Logger, d Deleter,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "handler.user.Delete")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(ctx, log, w, r)

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
