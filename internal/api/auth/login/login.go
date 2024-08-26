package login

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"providerHub/internal/api/auth/dto"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/logger/sl"
	"providerHub/pkg/validator"
	"time"
)

type UserProvider interface {
	Token(dto.LoginRequest) (token string, err error)
}

func New(
	log *slog.Logger,
	tokentTTL time.Duration,
	usrProvider UserProvider,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req dto.LoginRequest
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error parsing JSON", sl.Err(err))
			render.JSON(w, r, resp.Error("error to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, err)
			return
		}

		token, err := usrProvider.Token(req)
		if err != nil {
			log.Error("error during login", sl.Err(err))
			render.JSON(w, r, resp.Error("login failed"))
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
	}
}
