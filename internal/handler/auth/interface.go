package auth

import (
	"providerHub/internal/handler/auth/login"
	"providerHub/internal/handler/auth/refresh"
	"providerHub/internal/handler/auth/register"
)

type Interface interface {
	login.UserProvider
	refresh.UserRefresher
	register.UserRegister
}
