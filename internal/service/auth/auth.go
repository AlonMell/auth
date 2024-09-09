package auth

import (
	"context"
	"log/slog"
	"providerHub/internal/domain/dto"
	bc "providerHub/internal/infra/lib/bcrypt"
	"providerHub/internal/infra/lib/jwt"
	serInterface "providerHub/internal/service/interfaces"

	"providerHub/internal/domain/model"
	serr "providerHub/internal/service/errors"
)

//go:generate mockery --name Interface
type Interface interface {
	serInterface.UserSaver
	serInterface.UserEmailGetter
}

type Auth struct {
	log         *slog.Logger
	usrProvider Interface
}

func New(log *slog.Logger, p Interface) *Auth {
	return &Auth{log: log, usrProvider: p}
}

func (a *Auth) Token(
	ctx context.Context, tokenDTO dto.Token,
) (*dto.JWT, error) {
	const op = "service.Auth.Login"

	log := a.log.With(slog.String("op", op))
	log.Debug("login user")

	user, err := a.usrProvider.UserByEmail(ctx, tokenDTO.Email)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	//TODO: По ощущениям можно вызывать это всё в горутинах
	if err = bc.ComparePassword(user.PasswordHash, tokenDTO.Password); err != nil {
		return nil, serr.Catch(err, op)
	}

	access, err := jwt.GenerateToken(user.Id, user.Email, tokenDTO.AccessTTL, tokenDTO.Secret)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	refresh, err := jwt.GenerateToken(user.Id, user.Email, tokenDTO.RefreshTTL, tokenDTO.Secret)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	return &dto.JWT{Access: access, Refresh: refresh}, nil
}

func (a *Auth) RegisterUser(
	ctx context.Context, registerDTO dto.Register,
) (string, error) {
	const op = "service.Auth.Register"

	log := a.log.With(slog.String("op", op))
	log.Debug("registering user")

	hash, err := bc.GeneratePassword(registerDTO.Password)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	user := model.NewUser(registerDTO.Email, hash, true)

	id, err := a.usrProvider.SaveUser(ctx, *user)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return id, nil
}

func (a *Auth) RefreshToken(
	_ context.Context, refreshDTO dto.Refresh,
) (accessToken string, err error) {
	const op = "service.Auth.Refresh"

	log := a.log.With(slog.String("op", op))
	log.Debug("refresh token")

	claims, err := jwt.ValidateToken(refreshDTO.RefreshToken, refreshDTO.Secret)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	access, err := jwt.GenerateToken(claims.Subject, claims.Email, refreshDTO.AccessTTL, refreshDTO.Secret)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return access, nil
}
