package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"providerHub/internal/config"
	"providerHub/pkg/migrator"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config, logger *slog.Logger) (*Storage, error) {
	const op = "storage.postgres.New"

	sourceInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open("postgres", sourceInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	migratorConfig := migrator.New(db, cfg.SqlPath, cfg.Table, cfg.MajorVersion, cfg.MinorVersion)
	if err = migrator.Migrate(migratorConfig); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("Successfully connected to the database!")

	return &Storage{db}, nil
}
