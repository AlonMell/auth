package main

import "github.com/AlonMell/ProviderHub/internal/app/postgres"

type Config struct {
	Postgres  postgres.Config `yaml:"postgres" env-required:"true"`
	Migration `yaml:"migration" env-required:"true"`
}

type Migration struct {
	Path  string
	Table string
	Major int
	Minor int
}
