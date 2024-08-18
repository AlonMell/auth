package auth

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/router/handler/auth/dto"
	"providerHub/pkg/logger/sl"
)

type UserProvider interface {
	Token(dto.LoginRequest) (token string, err error)
}

func Login(log *slog.Logger, usrProvider UserProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "endpoint.auth.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req dto.LoginRequest
		req.Login = r.URL.Query().Get("login")
		req.Password = r.URL.Query().Get("password")

		log.Info("request query parameters", slog.Any("request", req))

		token, err := usrProvider.Token(req)
		if err != nil {
			log.Error("error during login", sl.Err(err))
			render.JSON(w, r, resp.Error("login failed"))
			return
		}

		w.Write([]byte(token))
	}
}
