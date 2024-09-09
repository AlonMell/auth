package user

import (
	"context"
	"log/slog"
	"providerHub/internal/domain/dto"
	bc "providerHub/internal/infra/lib/bcrypt"

	"providerHub/internal/domain/model"
	serr "providerHub/internal/service/errors"
	serInterface "providerHub/internal/service/interfaces"
)

type Interface interface {
	serInterface.UserSaver
	serInterface.UserIdGetter
	serInterface.UserUpdater
	serInterface.UserDeleter
}

type Provider struct {
	log         *slog.Logger
	usrProvider Interface
}

func New(log *slog.Logger, p Interface) *Provider {
	return &Provider{log: log, usrProvider: p}
}

func (p *Provider) Get(
	ctx context.Context, getDTO dto.UserGet,
) (u *model.User, err error) {
	const op = "user.Get"

	log := p.log.With(slog.String("op", op))
	log.Debug("get user from db")

	u, err = p.usrProvider.UserById(ctx, getDTO.Id)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	return
}

func (p *Provider) Create(
	ctx context.Context, createDTO dto.UserCreate,
) (string, error) {
	const op = "user.Create"

	log := p.log.With(slog.String("op", op))
	log.Debug("creating user")

	pass, err := bc.GeneratePassword(createDTO.Password)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	u := model.NewUser(createDTO.Email, pass, createDTO.IsActive)

	id, err := p.usrProvider.SaveUser(ctx, *u)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return id, nil
}

func (p *Provider) Delete(
	ctx context.Context, deleteDTO dto.UserDelete,
) error {
	const op = "user.Delete"

	log := p.log.With(slog.String("op", op))
	log.Debug("delete user from db")

	err := p.usrProvider.DeleteUser(ctx, deleteDTO.Id)
	if err != nil {
		return serr.Catch(err, op)
	}

	return nil
}

func (p *Provider) Update(
	ctx context.Context, updateDTO dto.UserUpdate,
) error {
	const op = "user.Update"

	log := p.log.With(slog.String("op", op))
	log.Debug("update user in db")

	pass, err := bc.GeneratePassword(updateDTO.Password)
	if err != nil {
		return serr.Catch(err, op)
	}

	u := model.NewUser(updateDTO.Email, pass, updateDTO.IsActive)

	err = p.usrProvider.UpdateUser(ctx, *u)
	if err != nil {
		return serr.Catch(err, op)
	}

	return nil
}
