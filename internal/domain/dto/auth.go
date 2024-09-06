package dto

import (
	"providerHub/internal/config"
)

type TokenDTO struct {
	Email    string
	Password string
	config.JWT
}

type RefreshDTO struct {
	RefreshToken string
	config.JWT
}

type RegisterDTO struct {
	Email    string
	Password string
}
