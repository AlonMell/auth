package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log/slog"
	"os"
	"providerHub/internal/config"
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

	if err = prepareDB(db, cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("Successfully connected to the database!")

	return &Storage{db}, nil
}

func prepareDB(db *sql.DB, cfg *config.Config) error {
	const op = "storage.postgres.prepareDB"

	file, err := os.Open(cfg.SqlPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	script, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(string(script))
	return err
}
