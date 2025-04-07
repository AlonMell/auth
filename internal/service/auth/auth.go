package auth

import (
	"context"
	"log/slog"

	"github.com/AlonMell/auth/internal/domain/dto"
	"github.com/AlonMell/auth/internal/domain/model"
	vo "github.com/AlonMell/auth/internal/domain/valueObject"
	bc "github.com/AlonMell/auth/internal/infra/lib/bcrypt"
	"github.com/AlonMell/auth/internal/infra/lib/jwt"
	"github.com/AlonMell/auth/internal/infra/lib/logger"
	ser "github.com/AlonMell/auth/internal/service"
	catcher "github.com/AlonMell/auth/internal/service/errors"
)

type UserRepo interface {
	ser.UserSaver
	ser.UserGetter
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

	user, err := a.usrProvider.User(ctx, vo.UserParams{"email": req.Email})
	if err != nil {
		return nil, catcher.Catch(ctx, err)
	}

	ctx = logger.WithLogUserID(ctx, user.Id)

	//TODO: По ощущениям можно вызывать это всё в горутинах
	if err = bc.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, catcher.Catch(ctx, err)
	}

	access, err := jwt.GenerateToken(user.Id, user.Email, a.jwt.AccessTTL, a.jwt.Secret)
	if err != nil {
		return nil, catcher.Catch(ctx, err)
	}

	refresh, err := jwt.GenerateToken(user.Id, user.Email, a.jwt.RefreshTTL, a.jwt.Secret)
	if err != nil {
		return nil, catcher.Catch(ctx, err)
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
		return "", catcher.Catch(ctx, err)
	}

	user := model.NewUser(req.Email, hash, true)

	ctx = logger.WithLogUserID(ctx, user.Id)

	id, err := a.usrProvider.SaveUser(ctx, *user)
	if err != nil {
		return "", catcher.Catch(ctx, err)
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
		return "", catcher.Catch(ctx, err)
	}

	access, err := jwt.GenerateToken(claims.Subject, claims.Email, a.jwt.AccessTTL, a.jwt.Secret)
	if err != nil {
		return "", catcher.Catch(ctx, err)
	}

	return access, nil
}
