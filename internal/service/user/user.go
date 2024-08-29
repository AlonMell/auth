package user

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"providerHub/internal/domain/model"
	"providerHub/internal/handler/user"
	bc "providerHub/internal/lib/bcrypt"
	"providerHub/internal/service"
	"providerHub/pkg/logger/sl"
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
	s   service.UserSaver
	g   Getter
	u   Updater
	d   Deleter
}

func New(
	log *slog.Logger,
	s service.UserSaver,
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

	u, err = p.g.UserById(r.UUID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (p *Provider) Create(r user.CreateUserRequest) (string, error) {
	const op = "user.Create"

	pass, err := bc.GeneratePassword(r.Password)
	if err != nil {
		p.log.Error("%v", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	u := model.User{
		UUID:         uuid.New().String(),
		Email:        r.Email,
		PasswordHash: pass,
		IsActive:     r.IsActive,
	}

	id, err := p.s.SaveUser(u)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *Provider) Delete(r user.DeleteUserRequest) error {
	const op = "user.Delete"

	err := p.d.DeleteUser(r.UUID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Provider) Update(r user.UpdateUserRequest) error {
	const op = "user.Update"

	pass, err := bc.GeneratePassword(r.Password)
	if err != nil {
		p.log.Error("%v", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	u := model.User{
		UUID:         r.UUID,
		Email:        r.Email,
		PasswordHash: pass,
		IsActive:     r.IsActive,
	}

	err = p.u.UpdateUser(u)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
