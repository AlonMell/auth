package config

import (
	"github.com/AlonMell/ProviderHub/internal/app/http"
	"github.com/AlonMell/ProviderHub/internal/app/postgres"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/jwt"
)

type Config struct {
	Env        string          `yaml:"env" env-required:"true"`
	Postgres   postgres.Config `yaml:"postgres" env-required:"true"`
	HTTPServer http.Config     `yaml:"http_server" env-required:"true"`
	JWT        jwt.Config      `yaml:"jwt" env-required:"true"`
}
