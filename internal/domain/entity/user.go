package entity

import (
	vo "github.com/AlonMell/auth/internal/domain/valueObject"
	"github.com/google/uuid"
)

type UserMap struct {
	Id string
	vo.UserParams
}

func NewUserMap(params map[string]any) *UserMap {
	return &UserMap{
		Id:         uuid.New().String(),
		UserParams: params,
	}
}
