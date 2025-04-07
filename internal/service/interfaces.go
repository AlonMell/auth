package service

import (
	"context"

	"github.com/AlonMell/auth/internal/domain/entity"
	"github.com/AlonMell/auth/internal/domain/model"
	vo "github.com/AlonMell/auth/internal/domain/valueObject"
)

type UserSaver interface {
	SaveUser(context.Context, model.User) (id string, err error)
}

type UserUpdater interface {
	UpdateUser(context.Context, entity.UserMap) error
}

type UserDeleter interface {
	DeleteUser(ctx context.Context, email string) error
}

type UserGetter interface {
	User(context.Context, vo.UserParams) (*model.User, error)
}
