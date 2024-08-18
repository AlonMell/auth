package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"providerHub/internal/domain/model"
	"providerHub/internal/lib/jwt"
	"providerHub/internal/router/handler/auth/dto"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
}

type UserSaver interface {
	SaveUser(model.User) (token string, err error)
}

type UserProvider interface {
	User(login string) (*model.User, error)
}

func New(
	log *slog.Logger,
	s UserSaver,
	p UserProvider,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    s,
		usrProvider: p,
	}
}

const secret = "secretsecretsecretsecretsecret"

func (a *Auth) Token(r dto.LoginRequest) (string, error) {
	const op = "service.Auth.Login"

	log := a.log.With(slog.String("op", op))

	log.Info("login user")

	user, err := a.usrProvider.User(r.Login)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(r.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	token, err := jwt.NewToken(*user, time.Hour, secret)
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err.Error())
	}

	return token, nil
}

func (a *Auth) RegisterUser(r dto.RegisterRequest) (string, error) {
	const op = "service.Auth.Register"

	log := a.log.With(slog.String("op", op))

	log.Info("registering user")

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err.Error())
	}

	user := model.User{
		Login:        r.Login,
		PasswordHash: hash,
	}

	userId, err := a.usrSaver.SaveUser(user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}
