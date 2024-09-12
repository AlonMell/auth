package auth

import (
	"github.com/AlonMell/ProviderHub/internal/handler/auth/login"
	"github.com/AlonMell/ProviderHub/internal/handler/auth/refresh"
	"github.com/AlonMell/ProviderHub/internal/handler/auth/register"
)

type Auth interface {
	login.UserProvider
	refresh.UserRefresher
	register.UserRegister
}
