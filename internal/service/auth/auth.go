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
	jwt         jwt.Config
	usrProvider UserRepo
}

func New(log *slog.Logger, p UserRepo, jwt jwt.Config) *Auth {
	return &Auth{log: log, usrProvider: p, jwt: jwt}
}

func (a *Auth) LoginUser(
	ctx context.Context, req dto.LoginReq,
) (*dto.JWT, error) {
	ctx = logger.WithLogOp(ctx, "service.Auth.Login")

	a.log.DebugContext(ctx, "login user")

	user, err := a.usrProvider.UserByEmail(ctx, req.Email)
	if err != nil {
		return nil, serr.Catch(ctx, err)
	}

	//TODO: По ощущениям можно вызывать это всё в горутинах
	if err = bc.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, serr.Catch(ctx, err)
	}

	access, err := jwt.GenerateToken(user.Id, user.Email, a.jwt.AccessTTL, a.jwt.Secret)
	if err != nil {
		return nil, serr.Catch(ctx, err)
	}

	refresh, err := jwt.GenerateToken(user.Id, user.Email, a.jwt.RefreshTTL, a.jwt.Secret)
	if err != nil {
		return nil, serr.Catch(ctx, err)
	}

	return &dto.JWT{Access: access, Refresh: refresh}, nil
}

func (a *Auth) RegisterUser(
	ctx context.Context, req dto.RegisterReq,
) (string, error) {
	ctx = logger.WithLogOp(ctx, "service.Auth.Register")

	a.log.DebugContext(ctx, "registering user")

	hash, err := bc.GeneratePassword(req.Password)
	if err != nil {
		return "", serr.Catch(ctx, err)
	}

	user := model.NewUser(req.Email, hash, true)
	ctx = logger.WithLogUserID(ctx, user.Id)

	id, err := a.usrProvider.SaveUser(ctx, *user)
	if err != nil {
		return "", serr.Catch(ctx, err)
	}

	return id, nil
}

func (a *Auth) RefreshToken(
	ctx context.Context, req dto.RefreshReq,
) (accessToken string, err error) {
	ctx = logger.WithLogOp(ctx, "service.Auth.Refresh")

	a.log.DebugContext(ctx, "refresh token")

	claims, err := jwt.ValidateToken(req.RefreshToken, a.jwt.Secret)
	if err != nil {
		return "", serr.Catch(ctx, err)
	}

	access, err := jwt.GenerateToken(claims.Subject, claims.Email, a.jwt.AccessTTL, a.jwt.Secret)
	if err != nil {
		return "", serr.Catch(ctx, err)
	}

	return access, nil
}
