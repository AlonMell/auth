package main

import (
	"flag"
	"fmt"
	"github.com/AlonMell/ProviderHub/cmd/migrator/config"
	"github.com/AlonMell/ProviderHub/internal/app/postgres"
	loader "github.com/AlonMell/ProviderHub/internal/infra/lib/config"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"github.com/AlonMell/migrator"
	_ "github.com/lib/pq"
)

func main() {
	var cfg config.Config
	loader.MustLoad(&cfg)
	storage, err := postgres.New(cfg.Postgres, logger.NewMockLogger())
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

	var path string
	flag.StringVar(&path, "path", "", "path to migrations")
	flag.Parse()
	if path == "" {
		panic("path to migrations is empty")
	}

	m := migrator.New(storage.DB, path, cfg.Table, cfg.Major, cfg.Minor)
	if err = m.Migrate(); err != nil {
		panic("failed to migrate: " + err.Error())
	}

	fmt.Println("migration completed successfully")
}
