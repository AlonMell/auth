package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Config struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Database string `yaml:"database" env-default:"postgres"`
}

// TODO: Добавить Пул потоков

func New(cfg Config, logger *slog.Logger) (*sqlx.DB, error) {
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
