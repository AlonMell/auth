package auth

import (
	"context"
	"github.com/AlonMell/ProviderHub/internal/domain/dto"
	"github.com/AlonMell/ProviderHub/internal/domain/model"
	bc "github.com/AlonMell/ProviderHub/internal/infra/lib/bcrypt"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/jwt"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	serr "github.com/AlonMell/ProviderHub/internal/service/errors"
	serInterface "github.com/AlonMell/ProviderHub/internal/service/interfaces"
	"log/slog"
)

type UserRepo interface {
	serInterface.UserSaver
	serInterface.UserEmailGetter
}

type Auth struct {
	log         *slog.Logger
	usrProvider UserRepo
}

func New(log *slog.Logger, p UserRepo) *Auth {
	return &Auth{log: log, usrProvider: p}
}

func (a *Auth) Token(
	ctx context.Context, tokenDTO dto.Token,
) (*dto.JWT, error) {
	const op = "service.Auth.Login"
	ctx = logger.WithLogOp(ctx, op)

	a.log.DebugContext(ctx, "login user")

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
	ctx = logger.WithLogOp(ctx, op)

	a.log.DebugContext(ctx, "registering user")

	hash, err := bc.GeneratePassword(registerDTO.Password)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	user := model.NewUser(registerDTO.Email, hash, true)
	ctx = logger.WithLogUserID(ctx, user.Id)

	id, err := a.usrProvider.SaveUser(ctx, *user)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return id, nil
}

func (a *Auth) RefreshToken(
	ctx context.Context, refreshDTO dto.Refresh,
) (accessToken string, err error) {
	const op = "service.Auth.Refresh"
	ctx = logger.WithLogOp(ctx, op)

	a.log.DebugContext(ctx, "refresh token")

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
