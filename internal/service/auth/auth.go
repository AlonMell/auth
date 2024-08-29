package auth

import (
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"providerHub/internal/handler/auth"
	bc "providerHub/internal/lib/bcrypt"
	"time"

	"providerHub/internal/domain/model"
	"providerHub/internal/lib/jwt"
	"providerHub/internal/service"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrGeneratingToken = errors.New("error generating token")
)

type UserGetter interface {
	UserByEmail(email string) (*model.User, error)
}

type Auth struct {
	log       *slog.Logger
	usrSaver  service.UserSaver
	usrGetter UserGetter
}

func New(
	log *slog.Logger,
	s service.UserSaver,
	g UserGetter,
) *Auth {
	return &Auth{
		log:       log,
		usrSaver:  s,
		usrGetter: g,
	}
}

const secret = "secretsecretsecretsecretsecret"

func (a *Auth) Token(r auth.LoginRequest, tokentTTL time.Duration) (string, error) {
	const op = "service.Auth.Login"

	log := a.log.With(slog.String("op", op))

	log.Info("login user")

	user, err := a.usrGetter.UserByEmail(r.Email)
	if err != nil {
		return "", service.Catch(err, op)
	}

	err = bc.ComparePassword(user.PasswordHash, r.Password)
	if err != nil {
		return "", service.Catch(err, op)
	}
	token, err := jwt.NewToken(*user, tokentTTL, secret)
	if err != nil {
		return "", service.Catch(err, op)
	}

	return token, nil
}

func (a *Auth) RegisterUser(r auth.RegisterRequest) (string, error) {
	const op = "service.Auth.Register"

	log := a.log.With(slog.String("op", op))

	log.Info("registering user")

	hash, err := bc.GeneratePassword(r.Password)
	if err != nil {
		return "", service.Catch(err, op)
	}

	user := model.User{
		UUID:         uuid.New().String(),
		Email:        r.Email,
		PasswordHash: hash,
	}

	id, err := a.usrSaver.SaveUser(user)
	if err != nil {
		return "", service.Catch(err, op)
	}

	return id, nil
}
