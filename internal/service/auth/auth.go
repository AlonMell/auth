package auth

import (
	"github.com/google/uuid"
	"log/slog"
	"providerHub/internal/config"
	"providerHub/internal/domain/dto"
	serInterface "providerHub/internal/service/interfaces"

	"providerHub/internal/domain/model"
	"providerHub/internal/handler/auth"
	bc "providerHub/internal/lib/bcrypt"
	"providerHub/internal/lib/jwt"
	serr "providerHub/internal/service/errors"
)

type UserGetter interface {
	UserByEmail(email string) (*model.User, error)
}

type Auth struct {
	log       *slog.Logger
	usrSaver  serInterface.UserSaver
	usrGetter UserGetter
}

func New(
	log *slog.Logger,
	s serInterface.UserSaver,
	g UserGetter,
) *Auth {
	return &Auth{
		log:       log,
		usrSaver:  s,
		usrGetter: g,
	}
}

func (a *Auth) Token(r auth.LoginRequest, cfg config.JWT) (*dto.JWT, error) {
	const op = "service.Auth.Login"

	log := a.log.With(slog.String("op", op))
	log.Debug("login user")

	user, err := a.usrGetter.UserByEmail(r.Email)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	if err = bc.ComparePassword(user.PasswordHash, r.Password); err != nil {
		return nil, serr.Catch(err, op)
	}

	access, err := jwt.GenerateToken(user.UUID, user.Email, cfg.AccessTTL, cfg.Secret)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	refresh, err := jwt.GenerateToken(user.UUID, user.Email, cfg.RefreshTTL, cfg.Secret)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	return &dto.JWT{Access: access, Refresh: refresh}, nil
}

func (a *Auth) RegisterUser(r auth.RegisterRequest) (string, error) {
	const op = "service.Auth.Register"

	log := a.log.With(slog.String("op", op))
	log.Debug("registering user")

	hash, err := bc.GeneratePassword(r.Password)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	user := model.User{
		UUID:         uuid.New().String(),
		Email:        r.Email,
		PasswordHash: hash,
	}

	id, err := a.usrSaver.SaveUser(user)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return id, nil
}

func (a *Auth) RefreshToken(req auth.RefreshRequest, cfg config.JWT) (accessToken string, err error) {
	const op = "service.Auth.Refresh"

	log := a.log.With(slog.String("op", op))
	log.Debug("refresh token")

	claims, err := jwt.ValidateToken(req.RefreshToken, cfg.Secret)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	access, err := jwt.GenerateToken(claims.Subject, claims.Email, cfg.AccessTTL, cfg.Secret)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return access, nil
}
