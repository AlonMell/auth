package register

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
)

type UserRegister interface {
	RegisterUser(dto.RegisterRequest) (uuid string, err error)
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

		// Системная ошибка / Пользовательская ошибка
		err := decoder.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("error parsing JSON", sl.Err(err))
			render.JSON(w, r, resp.Error("error to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Пользовательская ошибка
		if err = validator.Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, err)
			return
		}

		// Программная ошибка / Пользовательская ошибка
		uuid, err := reg.RegisterUser(req)
		if err != nil {
			render.JSON(w, r, resp.Error("error to register user"))
			return
		}

		render.JSON(w, r, uuid)
	}
}
