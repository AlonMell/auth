package get

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/domain/model"
	"providerHub/internal/handler"
	"providerHub/internal/handler/user"
	resp "providerHub/internal/lib/api/response"
	"providerHub/pkg/validator"
)

type Getter interface {
	Get(user.GetUserRequest) (*model.User, error)
}

func New(log *slog.Logger, g Getter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.get.New"

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req user.GetUserRequest
		req.UUID = r.URL.Query().Get("id")

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		u, err := g.Get(req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		resp.WriteJSON(
			w, r,
			user.GetUserResponse{
				Email:        u.Email,
				PasswordHash: u.PasswordHash,
				IsActive:     u.IsActive,
			},
		)
	}
}
