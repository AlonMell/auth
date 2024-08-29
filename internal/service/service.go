package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"providerHub/internal/domain/model"
	repo "providerHub/internal/repository"
	"providerHub/internal/service/auth"
)

type UserSaver interface {
	SaveUser(model.User) (uuid string, err error)
}

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
		return New(auth.ErrInvalidPassword, UserKind)
	case errors.Is(err, auth.ErrGeneratingToken):
		return New(err, InternalKind)
	default:
		return New(fmt.Errorf("%s: %w", op, err), InternalKind)
	}
}
