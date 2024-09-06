package register

import (
	"context"
	"log/slog"
	"net/http"
	"providerHub/internal/domain/dto"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/handler"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/validator"
)

type UserRegister interface {
	RegisterUser(context.Context, dto.RegisterDTO) (id string, err error)
}

// New
// @Summary Register
// @Tags auth
// @Description Register new user
// @Accept json
// @Produce json
// @Param input body Request true "user info"
// @Success 200 {object} Response
// @Router /register [post]
func New(
	log *slog.Logger, reg UserRegister,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.register.New"

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := decoder.DecodeJSON(r.Body, &req)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		registerDTO := dto.RegisterDTO{
			Email:    req.Email,
			Password: req.Password,
		}

		id, err := reg.RegisterUser(r.Context(), registerDTO)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{Id: id})
	}
}
