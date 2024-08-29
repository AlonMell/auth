package post

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/handler/user"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/logger/sl"
	"providerHub/pkg/validator"
)

type Creater interface {
	Create(user.CreateUserRequest) (string, error)
}

func New(log *slog.Logger, c Creater) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.post.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req user.CreateUserRequest
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error parsing JSON", sl.Err(err))
			http.Error(w, "error to decode request", http.StatusBadRequest)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		uuid, err := c.Create(req)
		if err != nil {
			log.Error("error during create user", sl.Err(err))
			http.Error(w, "error to create user", http.StatusInternalServerError)
			return
		}

		resp.WriteJSON(w, r, Ok(uuid))
	}
}

func Ok(uuid string) user.CreateUserResponse {
	return user.CreateUserResponse{
		UUID:     uuid,
		Response: resp.Ok(),
	}
}
