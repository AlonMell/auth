package register

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/handler/auth"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/logger/sl"
	"providerHub/pkg/validator"
)

type UserRegister interface {
	RegisterUser(auth.RegisterRequest) (uuid string, err error)
}

func New(
	log *slog.Logger,
	reg UserRegister,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.register.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req auth.RegisterRequest
		err := decoder.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("error parsing JSON", sl.Err(err))
			resp.WriteJSON(w, r, resp.Error("error to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			resp.WriteJSON(w, r, err)
			return
		}

		uuid, err := reg.RegisterUser(req)
		if err != nil {
			resp.WriteJSON(w, r, resp.Error("error to register user"))
			return
		}

		resp.WriteJSON(w, r, Ok(uuid))
	}
}

func Ok(uuid string) auth.RegisterResponse {
	return auth.RegisterResponse{
		UUID:     uuid,
		Response: resp.Ok(),
	}
}
