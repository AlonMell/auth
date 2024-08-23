package register

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/internal/router/handler/auth/dto"
	"providerHub/pkg/logger/sl"
	"providerHub/pkg/validator"
)

type UserRegister interface {
	RegisterUser(dto.RegisterRequest) (userId string, err error)
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

		var req dto.RegisterRequest
		err := decoder.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("error parsing JSON", sl.Err(err))
			render.JSON(w, r, resp.Error("error to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, err)
			return
		}

		userId, err := reg.RegisterUser(req)
		if err != nil {
			render.JSON(w, r, resp.Error("error to register user"))
			return
		}

		render.JSON(w, r, userId)
	}
}
