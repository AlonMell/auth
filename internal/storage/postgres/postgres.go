package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"providerHub/internal/config"
	"providerHub/internal/domain/models"
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
		return nil, err
	}

	logger.Info("Successfully connected to the database!")

	return &Storage{db}, nil
}

func (s *Storage) SaveUser(user models.User) error {
	const op = "storage.postgres.SaveUser"

	query := `INSERT INTO users(login, email, password_hash, phone, is_active) VALUES ($1, $2, $3, $4, $5)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Login, user.Email, user.PasswordHash, user.Phone, user.IsActive)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserByLogin(login string) (*models.User, error) {
	const op = "storage.postgres.GetUserByLogin"

	query := `
		SELECT login, password_hash, phone, email, is_acrive 
		FROM users 
		WHERE login=$1`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(login).Scan(&user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
