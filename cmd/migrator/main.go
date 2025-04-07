package main

import (
	"fmt"

	"github.com/AlonMell/auth/cmd/migrator/config"
	"github.com/AlonMell/auth/internal/app/postgres"
	loader "github.com/AlonMell/auth/internal/infra/lib/config"
	"github.com/AlonMell/grovelog/util"
	"github.com/AlonMell/migrator"
	_ "github.com/lib/pq"
)

func main() {
	paths := config.MustLoadFlags()

	var cfg config.Config
	loader.MustLoad(&cfg, paths.Config)

	storage, err := postgres.New(cfg.Postgres, util.NewMockLogger())
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

	m := migrator.New(storage.DB, paths.Migrations, cfg.Table, cfg.Major, cfg.Minor)
	if err = m.Migrate(); err != nil {
		panic("failed to migrate: " + err.Error())
	}

	fmt.Println("migration completed successfully")
}
