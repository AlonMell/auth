package user

import (
	"context"
	"github.com/AlonMell/ProviderHub/internal/domain/dto"
	bc "github.com/AlonMell/ProviderHub/internal/infra/lib/bcrypt"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"log/slog"

	"github.com/AlonMell/ProviderHub/internal/domain/model"
	serr "github.com/AlonMell/ProviderHub/internal/service/errors"
	serInterface "github.com/AlonMell/ProviderHub/internal/service/interfaces"
)

type UserRepo interface {
	serInterface.UserSaver
	serInterface.UserIdGetter
	serInterface.UserUpdater
	serInterface.UserDeleter
}

type Provider struct {
	log         *slog.Logger
	usrProvider UserRepo
}

func New(log *slog.Logger, p UserRepo) *Provider {
	return &Provider{log: log, usrProvider: p}
}

func (p *Provider) Get(
	ctx context.Context, getDTO dto.UserGet,
) (*model.User, error) {
	const op = "service.user.Get"
	ctx = logger.WithLogOp(ctx, op)

	p.log.DebugContext(ctx, "get user from db")

	user, err := p.usrProvider.UserById(ctx, getDTO.Id)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	ctx = logger.WithLogUserID(ctx, user.Id)

	return user, err
}

func (p *Provider) Create(
	ctx context.Context, createDTO dto.UserCreate,
) (string, error) {
	const op = "service.user.Create"
	ctx = logger.WithLogOp(ctx, op)

	p.log.DebugContext(ctx, "creating user")

	pass, err := bc.GeneratePassword(createDTO.Password)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	u := model.NewUser(createDTO.Email, pass, createDTO.IsActive)

	ctx = logger.WithLogUserID(ctx, u.Id)

	id, err := p.usrProvider.SaveUser(ctx, *u)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return id, nil
}

func (p *Provider) Delete(
	ctx context.Context, deleteDTO dto.UserDelete,
) error {
	const op = "service.user.Delete"
	ctx = logger.WithLogOp(ctx, op)

	p.log.DebugContext(ctx, "delete user from db")

	err := p.usrProvider.DeleteUser(ctx, deleteDTO.Id)
	if err != nil {
		return serr.Catch(err, op)
	}

	return nil
}

func (p *Provider) Update(
	ctx context.Context, updateDTO dto.UserUpdate,
) error {
	const op = "service.user.Update"
	ctx = logger.WithLogOp(ctx, op)

	p.log.DebugContext(ctx, "update user in db")

	pass, err := bc.GeneratePassword(updateDTO.Password)
	if err != nil {
		return serr.Catch(err, op)
	}

	u := model.NewUser(updateDTO.Email, pass, updateDTO.IsActive)

	ctx = logger.WithLogUserID(ctx, u.Id)

	err = p.usrProvider.UpdateUser(ctx, *u)
	if err != nil {
		return serr.Catch(err, op)
	}

	return nil
}
