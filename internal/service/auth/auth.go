package auth

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"providerHub/internal/handler/auth"
	"time"

	"providerHub/internal/domain/model"
	"providerHub/internal/lib/jwt"
	"providerHub/internal/service"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
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

func (a *Auth) Token(r auth.LoginRequest) (string, error) {
	const op = "service.Auth.Login"

	log := a.log.With(slog.String("op", op))

	log.Info("login user")

	user, err := a.usrGetter.UserByEmail(r.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(r.Password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidPassword)
	}
	token, err := jwt.NewToken(*user, time.Hour, secret)
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	return token, nil
}

func (a *Auth) RegisterUser(r auth.RegisterRequest) (string, error) {
	const op = "service.Auth.Register"

	log := a.log.With(slog.String("op", op))

	log.Info("registering user")

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}

	user := model.User{
		UUID:         uuid.New().String(),
		Email:        r.Email,
		PasswordHash: hash,
	}

	id, err := a.usrSaver.SaveUser(user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
