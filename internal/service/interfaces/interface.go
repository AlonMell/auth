package interfaces

import (
	"context"
	"providerHub/internal/domain/model"
)

type UserSaver interface {
	SaveUser(context.Context, model.User) (id string, err error)
}

type UserUpdater interface {
	UpdateUser(context.Context, model.User) error
}

type UserDeleter interface {
	DeleteUser(ctx context.Context, email string) error
}

type UserGetter interface {
	UserById(context.Context, string) (*model.User, error)
	UserByEmail(context.Context, string) (*model.User, error)
}
