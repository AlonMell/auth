package delete

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

type Deleter interface {
	Delete(user.DeleteUserRequest) error
}

func New(log *slog.Logger, d Deleter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req user.DeleteUserRequest
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error parsing JSON", sl.Err(err))
			http.Error(w, "error to decode request", http.StatusBadRequest)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			resp.WriteJSON(w, r, err)
			return
		}

		if err := d.Delete(req); err != nil {
			log.Error("error during delete user", sl.Err(err))
			http.Error(w, "error to delete user", http.StatusInternalServerError)
			return
		}

		resp.WriteJSON(w, r, Ok())
	}
}

func Ok() user.DeleteUserResponse {
	return user.DeleteUserResponse{
		Response: resp.Ok(),
	}
}
