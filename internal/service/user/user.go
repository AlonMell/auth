package user

import (
	"log/slog"

	"github.com/google/uuid"

	"providerHub/internal/domain/model"
	"providerHub/internal/handler/user"
	bc "providerHub/internal/lib/bcrypt"
	serr "providerHub/internal/service/errors"
	serInterface "providerHub/internal/service/interfaces"
)

type Updater interface {
	UpdateUser(user model.User) error
}

type Deleter interface {
	DeleteUser(email string) error
}

type Getter interface {
	UserById(email string) (*model.User, error)
}

type Provider struct {
	log *slog.Logger
	s   serInterface.UserSaver
	g   Getter
	u   Updater
	d   Deleter
}

func New(
	log *slog.Logger,
	s serInterface.UserSaver,
	g Getter,
	u Updater,
	d Deleter,
) *Provider {
	return &Provider{
		log: log,
		s:   s,
		g:   g,
		u:   u,
		d:   d,
	}
}

func (p *Provider) Get(r user.GetUserRequest) (u *model.User, err error) {
	const op = "user.Get"

	log := p.log.With(slog.String("op", op))
	log.Debug("get user from db")

	u, err = p.g.UserById(r.UUID)
	if err != nil {
		return nil, serr.Catch(err, op)
	}

	return
}

func (p *Provider) Create(r user.CreateUserRequest) (string, error) {
	const op = "user.Create"

	log := p.log.With(slog.String("op", op))
	log.Debug("creating user")

	pass, err := bc.GeneratePassword(r.Password)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	u := model.User{
		UUID:         uuid.New().String(),
		Email:        r.Email,
		PasswordHash: pass,
		IsActive:     r.IsActive,
	}

	id, err := p.s.SaveUser(u)
	if err != nil {
		return "", serr.Catch(err, op)
	}

	return id, nil
}

func (p *Provider) Delete(r user.DeleteUserRequest) error {
	const op = "user.Delete"

	log := p.log.With(slog.String("op", op))
	log.Debug("delete user from db")

	err := p.d.DeleteUser(r.UUID)
	if err != nil {
		return serr.Catch(err, op)
	}

	return nil
}

func (p *Provider) Update(r user.UpdateUserRequest) error {
	const op = "user.Update"

	log := p.log.With(slog.String("op", op))
	log.Debug("update user in db")

	pass, err := bc.GeneratePassword(r.Password)
	if err != nil {
		return serr.Catch(err, op)
	}

	u := model.User{
		UUID:         r.UUID,
		Email:        r.Email,
		PasswordHash: pass,
		IsActive:     r.IsActive,
	}

	err = p.u.UpdateUser(u)
	if err != nil {
		return serr.Catch(err, op)
	}

	return nil
}
