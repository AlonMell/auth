package get

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"providerHub/internal/domain/model"
	"providerHub/internal/handler/user"
	resp "providerHub/internal/lib/api/response"
	"providerHub/pkg/logger/sl"
	"providerHub/pkg/validator"
)

type Getter interface {
	Get(user.GetUserRequest) (*model.User, error)
}

func New(log *slog.Logger, g Getter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req user.GetUserRequest
		//TODO: refactor this
		req.UUID = r.URL.Query().Get("id")

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			resp.WriteJSON(w, r, err)
			return
		}

		u, err := g.Get(req)
		if err != nil {
			log.Error("error during get user", sl.Err(err))
			resp.WriteJSON(w, r, resp.Error("get user failed"))
			return
		}

		resp.WriteJSON(w, r, Ok(u))
	}
}

func Ok(u *model.User) user.GetUserResponse {
	return user.GetUserResponse{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		IsActive:     u.IsActive,
		Response:     resp.Ok(),
	}
}
