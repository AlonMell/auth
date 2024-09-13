package auth

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

type UserRegister interface {
	RegisterUser(context.Context, dto.RegisterReq) (id string, err error)
}

type RegisterResp struct {
	Id string `json:"id"`
	resp.Response
}

// Register
// @Summary Register
// @Tags auth
// @Description Register new user
// @Accept json
// @Produce json
// @Param input body Request true "user info"
// @Success 200 {object} Response
// @Router /register [post]
func Register(
	log *slog.Logger, reg UserRegister,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "handler.auth.Register")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(ctx, log, w, r)

		var req dto.RegisterReq
		err := decoder.DecodeJSON(r.Body, &req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		log.InfoContext(ctx, "request body decoded", slog.Any("request", req))

		if err = validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		id, err := reg.RegisterUser(ctx, req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.Status(r, http.StatusOK)
		resp.WriteJSON(w, r, RegisterResp{Id: id})
	}
}
