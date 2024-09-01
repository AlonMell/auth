package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	_ "github.com/lib/pq"

	"providerHub/internal/config"
	"providerHub/internal/domain/model"
	repo "providerHub/internal/repository"
)

// TODO: Добавить Пул потоков

type Storage struct {
	db *sql.DB
	mu sync.RWMutex
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

	logger.Info("Successfully connected to the database!")

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(user model.User) (string, error) {
	const op = "storage.postgres.SaveUser"
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.userById(user.UUID)
	if err == nil {
		return "", repo.ErrUserExists
	} else if !errors.Is(err, repo.ErrUserNotFound) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	query := `INSERT INTO users(id, email, password_hash, is_active) VALUES ($1, $2, $3, $4)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UUID, user.Email, user.PasswordHash /*user.Phone,*/, user.IsActive)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	usr, err := s.userById(user.UUID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return usr.UUID, nil
}

func (s *Storage) UserById(uuid string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userById(uuid)
}

func (s *Storage) userById(uuid string) (*model.User, error) {
	const op = "storage.postgres.UserById"

	query := `
		SELECT id, email, password_hash, is_active 
		FROM users 
		WHERE id=$1`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(uuid).Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.IsActive)
	if err != nil {
		return nil, errorHandler(err, op)
	}

	return &user, nil
}

func (s *Storage) UserByEmail(email string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userByEmail(email)
}

func (s *Storage) DeleteUser(uuid string) error {
	const op = "storage.postgres.DeleteUser"
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM users WHERE id=$1`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateUser(user model.User) error {
	const op = "storage.postgres.UpdateUser"
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `UPDATE users SET email=$1, password_hash=$2, is_active=$3 WHERE id=$4`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.PasswordHash, user.IsActive, user.UUID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) userByEmail(email string) (*model.User, error) {
	const op = "storage.postgres.User"

	query := `
		SELECT id, email, password_hash, is_active 
		FROM users 
		WHERE email=$1`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(email).Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.IsActive)
	if err != nil {
		return nil, errorHandler(err, op)
	}

	return &user, nil
}

func errorHandler(err error, op string) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return repo.ErrUserNotFound
	default:
		return fmt.Errorf("%s: %w", op, err)
	}
}
