package service

import "providerHub/internal/domain/model"

type UserSaver interface {
	SaveUser(model.User) (uuid string, err error)
}
