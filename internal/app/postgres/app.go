package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"providerHub/internal/config"
)

// TODO: Добавить Пул потоков

func New(cfg *config.Config, logger *slog.Logger) (*sql.DB, error) {
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

	logger.Info("Successfully connected to the database!")

	return db, nil
}
