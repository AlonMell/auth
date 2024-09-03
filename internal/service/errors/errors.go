package errors

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"providerHub/internal/lib/jwt"
	repo "providerHub/internal/repo"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

const (
	UserKind = iota + 1
	InternalKind
	SystemKind
)

type CustomError struct {
	Err  error
	Kind int
}

func (c *CustomError) Error() string {
	return c.Unwrap().Error()
}

func (c *CustomError) Unwrap() error {
	return c.Err
}

func New(err error, kind int) *CustomError {
	return &CustomError{
		Err:  err,
		Kind: kind,
	}
}

func Catch(err error, op string) error {
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return New(err, UserKind)
	case errors.Is(err, repo.ErrUserExists):
		return New(err, UserKind)
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return New(ErrInvalidPassword, UserKind)
	case errors.Is(err, jwt.ErrGeneratingToken):
		return New(err, InternalKind)
	default:
		return New(fmt.Errorf("%s: %w", op, err), InternalKind)
	}
}
