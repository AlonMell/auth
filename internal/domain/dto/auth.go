package dto

import (
	"providerHub/internal/infra/config"
)

type Token struct {
	Email    string
	Password string
	config.JWT
}

type Refresh struct {
	RefreshToken string
	config.JWT
}

type Register struct {
	Email    string
	Password string
}
