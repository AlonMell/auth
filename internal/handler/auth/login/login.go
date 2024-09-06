package login

import (
	"context"
	"log/slog"
	"net/http"
	"providerHub/internal/config"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"providerHub/internal/domain/dto"
	"providerHub/internal/handler"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/lib/decoder"
	"providerHub/pkg/validator"
)

type UserProvider interface {
	Token(context.Context, dto.TokenDTO) (jwt *dto.JWT, err error)
}

// New
// @Summary Sign In
// @Tags auth
// @Description Sign in to the system
// @Accept json
// @Produce json
// @Param input body Request true "account info"
// @Success 200 {object} Response
// @Router /login [post]
func New(
	log *slog.Logger, cfg config.JWT, usrProvider UserProvider,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.login.New"

		errCatcher := handler.NewCatcher(op, log, w, r)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := decoder.DecodeJSON(r.Body, &req); err != nil {
			errCatcher.Catch(err)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.Struct(req); err != nil {
			errCatcher.Catch(err)
			return
		}

		tokenDTO := dto.TokenDTO{
			Email:    req.Email,
			Password: req.Password,
			JWT:      cfg,
		}

		jwt, err := usrProvider.Token(r.Context(), tokenDTO)
		if err != nil {
			errCatcher.Catch(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "refresh",
			Value: jwt.Refresh,
			//Path:     "/",
			HttpOnly: true,
			Secure:   false,
			//SameSite: http.SameSiteLaxMode,
			Expires: time.Now().Add(cfg.RefreshTTL),
			//Domain:   "example.com",
		})

		w.WriteHeader(http.StatusOK)
		resp.WriteJSON(w, r, Response{Jwt: jwt.Access})
	}
}
