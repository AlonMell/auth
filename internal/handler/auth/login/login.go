package login

import (
	"log/slog"
	"net/http"
	"providerHub/internal/config"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/domain/dto"
	"providerHub/internal/handler"
	"providerHub/internal/handler/auth"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/validator"
)

type UserProvider interface {
	Token(auth.LoginRequest, config.JWT) (jwt *dto.JWT, err error)
}

func New(
	log *slog.Logger,
	cfg config.JWT,
	usrProvider UserProvider,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.login.New"

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req auth.LoginRequest
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		jwt, err := usrProvider.Token(req, cfg)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		res := auth.LoginResponse{Jwt: jwt.Access}

		http.SetCookie(w, &http.Cookie{
			Name:  "refresh",
			Value: jwt.Refresh,
			//Path:     "/",
			HttpOnly: true,
			Secure:   false,
			//SameSite: http.SameSiteLaxMode,
			Expires: time.Now().Add(cfg.RefreshTTL),
			//Domain:   "example.com",
		})

		resp.WriteJSON(w, r, res)
	}
}
