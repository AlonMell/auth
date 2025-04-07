package auth

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/AlonMell/auth/internal/handler/errors"
	resp "github.com/AlonMell/auth/internal/infra/lib/api/response"
	"github.com/AlonMell/auth/internal/infra/lib/decoder"
	"github.com/AlonMell/auth/internal/infra/lib/logger"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlonMell/auth/internal/domain/dto"
	"github.com/AlonMell/auth/pkg/validator"
)

type UserProvider interface {
	LoginUser(context.Context, dto.LoginReq) (jwt *dto.JWT, err error)
}

type LoginResp struct {
	Jwt string `json:"jwt"`
}

func Login(
	log *slog.Logger, refreshTTL time.Duration, usrProvider UserProvider,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogOp(r.Context(), "handler.auth.Login")
		ctx = logger.WithLogRequestID(ctx, middleware.GetReqID(ctx))

		errCatcher := errors.NewCatcher(ctx, log, w, r)

		var req dto.LoginReq
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.InfoContext(ctx, "request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		jwt, err := usrProvider.LoginUser(ctx, req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "refresh",
			Value: jwt.Refresh,
			//Path:     "/",
			HttpOnly: true,
			Secure:   false,
			//SameSite: http.SameSiteLaxMode,
			Expires: time.Now().Add(refreshTTL),
			//Domain:   "example.com",
		})

		resp.Status(r, http.StatusOK)
		resp.WriteJSON(w, r, LoginResp{Jwt: jwt.Access})
	}
}
