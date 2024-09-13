package errors

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/jwt"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"github.com/AlonMell/ProviderHub/internal/infra/repo"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

const (
	BadRequest   = http.StatusBadRequest
	Unauthorized = http.StatusUnauthorized
	NotFound     = http.StatusNotFound
	Internal     = http.StatusInternalServerError
)

const (
	UserKind = iota + 1
	InternalKind
	SystemKind
)

type ServiceError struct {
	Err  error
	Kind int
	Code int
}

func (c *ServiceError) Error() string {
	return c.Err.Error()
}

func (c *ServiceError) Unwrap() error {
	return c.Err
}

func Wrap(err error, kind int, code int) *ServiceError {
	return &ServiceError{
		Err:  err,
		Kind: kind,
		Code: code,
	}
}

func WrapCtx(ctx context.Context, err error) error {
	return logger.Wrap(ctx, err)
}

func Catch(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return WrapCtx(ctx, Wrap(err, UserKind, NotFound))
	case errors.Is(err, repo.ErrUserExists):
		fmt.Println("Зашел сюда." + err.Error())
		return WrapCtx(ctx, Wrap(err, UserKind, BadRequest))
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return WrapCtx(ctx, Wrap(ErrInvalidPassword, UserKind, BadRequest))
	case errors.Is(err, jwt.ErrGeneratingToken):
		return WrapCtx(ctx, Wrap(err, InternalKind, Internal))
	case errors.Is(err, jwt.ErrValidatingToken):
		return WrapCtx(ctx, Wrap(err, UserKind, Unauthorized))
	default:
		return WrapCtx(ctx, Wrap(err, InternalKind, Internal))
	}
}
