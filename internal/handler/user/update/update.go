package update

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/handler"
	"providerHub/internal/handler/user"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/validator"
)

type Updater interface {
	Update(user.UpdateUserRequest) error
}

func New(log *slog.Logger, u Updater) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.update.New"

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req user.UpdateUserRequest
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		if err := u.Update(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.WriteJSON(w, r, Ok())
	}
}

func Ok() user.UpdateUserResponse {
	return user.UpdateUserResponse{Response: resp.Ok()}
}
