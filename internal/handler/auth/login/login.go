package login

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/handler"
	"providerHub/internal/handler/auth"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/validator"
	"time"
)

type UserProvider interface {
	Token(auth.LoginRequest, time.Duration) (token string, err error)
}

func New(
	log *slog.Logger,
	tokentTTL time.Duration,
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

		token, err := usrProvider.Token(req, tokentTTL)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: token,
			//Path:     "/",
			HttpOnly: true,
			Secure:   false,
			//SameSite: http.SameSiteLaxMode,
			Expires: time.Now().Add(tokentTTL),
			//Domain:   "example.com",
		})

		resp.WriteJSON(w, r, resp.Ok())
	}
}
