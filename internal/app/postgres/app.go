package postgres

import (
	"fmt"
	"github.com/AlonMell/ProviderHub/internal/infra/config"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

// TODO: Добавить Пул потоков

func New(cfg *config.Config, logger *slog.Logger) (*sqlx.DB, error) {
	const op = "storage.postgres.New"

	sourceInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sqlx.Connect("postgres", sourceInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("Successfully connected to the database!")

	return db, nil
}
