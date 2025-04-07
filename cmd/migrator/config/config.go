package config

import "github.com/AlonMell/auth/internal/app/postgres"

type Config struct {
	Postgres  postgres.Config `yaml:"postgres" env-required:"true"`
	Migration `yaml:"migration" env-required:"true"`
}

type Migration struct {
	Table string
	Major int
	Minor int
}
