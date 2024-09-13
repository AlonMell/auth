package main

import (
	"fmt"
	"github.com/AlonMell/ProviderHub/internal/app/postgres"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/config"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"github.com/AlonMell/migrator"
	_ "github.com/lib/pq"
)

func main() {
	var cfg Config
	config.MustLoad(&cfg)
	storage, err := postgres.New(cfg.Postgres, logger.NewMockLogger())
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

	m := migrator.New(storage.DB, cfg.Path, cfg.Table, cfg.Major, cfg.Minor)
	if err = m.Migrate(); err != nil {
		panic("failed to migrate: " + err.Error())
	}

	fmt.Println("migration completed successfully")
}
