package user

import (
	"context"
	"log/slog"

	"github.com/AlonMell/auth/internal/domain/dto"
	"github.com/AlonMell/auth/internal/domain/entity"
	vo "github.com/AlonMell/auth/internal/domain/valueObject"
	bc "github.com/AlonMell/auth/internal/infra/lib/bcrypt"
	"github.com/AlonMell/auth/internal/infra/lib/logger"
	ser "github.com/AlonMell/auth/internal/service"
	catcher "github.com/AlonMell/auth/internal/service/errors"

	"github.com/AlonMell/auth/internal/domain/model"
)

type Repo interface {
	ser.UserSaver
	ser.UserGetter
	ser.UserUpdater
	ser.UserDeleter
}

type Provider struct {
	log         *slog.Logger
	usrProvider Repo
}

func New(log *slog.Logger, p Repo) *Provider {
	return &Provider{log: log, usrProvider: p}
}

func (p *Provider) Get(
	ctx context.Context, req dto.UserGetReq,
) (*model.User, error) {
	ctx = logger.WithLogOp(ctx, "service.user.Get")

	p.log.DebugContext(ctx, "get user from db")

	user, err := p.usrProvider.User(ctx, vo.UserParams{"id": req.Id})
	if err != nil {
		return nil, catcher.Catch(ctx, err)
	}

	ctx = logger.WithLogUserID(ctx, user.Id)

	return user, err
}

func (p *Provider) Create(
	ctx context.Context, req dto.UserCreateReq,
) (string, error) {
	ctx = logger.WithLogOp(ctx, "service.user.Create")

	p.log.DebugContext(ctx, "creating user")

	pass, err := bc.GeneratePassword(req.Password)
	if err != nil {
		return "", catcher.Catch(ctx, err)
	}

	u := model.NewUser(req.Email, pass, req.IsActive)

	ctx = logger.WithLogUserID(ctx, u.Id)

	id, err := p.usrProvider.SaveUser(ctx, *u)
	if err != nil {
		return "", catcher.Catch(ctx, err)
	}

	return id, nil
}

func (p *Provider) Delete(
	ctx context.Context, req dto.UserDeleteReq,
) error {
	ctx = logger.WithLogOp(ctx, "service.user.Delete")

	p.log.DebugContext(ctx, "delete user from db")

	err := p.usrProvider.DeleteUser(ctx, req.Id)
	if err != nil {
		return catcher.Catch(ctx, err)
	}

	return nil
}

func (p *Provider) Update(
	ctx context.Context, req dto.UserUpdateReq,
) error {
	ctx = logger.WithLogOp(ctx, "service.user.Update")

	p.log.DebugContext(ctx, "update user in db")

	pass, err := bc.GeneratePassword(req.Password)
	if err != nil {
		return catcher.Catch(ctx, err)
	}

	u := entity.NewUserMap(map[string]any{
		"email":     req.Email,
		"password":  pass,
		"is_active": req.IsActive,
	})

	ctx = logger.WithLogUserID(ctx, u.Id)

	err = p.usrProvider.UpdateUser(ctx, *u)
	if err != nil {
		return catcher.Catch(ctx, err)
	}

	return nil
}
