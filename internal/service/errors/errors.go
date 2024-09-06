package errors

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"providerHub/internal/lib/jwt"
	repo "providerHub/internal/repo"
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

type CustomError struct {
	Err  error
	Kind int
	Code int
}

func (c *CustomError) Error() string {
	return c.Unwrap().Error()
}

func (c *CustomError) Unwrap() error {
	return c.Err
}

func New(err error, kind int, code int) *CustomError {
	return &CustomError{
		Err:  err,
		Kind: kind,
		Code: code,
	}
}

func Catch(err error, op string) error {
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return New(err, UserKind, NotFound)
	case errors.Is(err, repo.ErrUserExists):
		return New(err, UserKind, BadRequest)
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return New(ErrInvalidPassword, UserKind, BadRequest)
	case errors.Is(err, jwt.ErrGeneratingToken):
		return New(err, InternalKind, Internal)
	case errors.Is(err, jwt.ErrValidatingToken):
		return New(err, UserKind, Unauthorized)
	default:
		return New(fmt.Errorf("%s: %w", op, err), InternalKind, Internal)
	}
}
