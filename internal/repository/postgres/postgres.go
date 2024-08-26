package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"providerHub/internal/config"
	"providerHub/internal/domain/model"
)

const (
	ErrUserExists = "user already exists"
)

// TODO: Добавить RWMutex и Пул потоков
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

	logger.Info("Successfully connected to the database!")

	return &Storage{db}, nil
}

func (s *Storage) SaveUser(user model.User) (uuid string, err error) {
	const op = "storage.postgres.SaveUser"

	exists, err := s.IsUserExists(user.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	} else if exists {
		return "", errors.New(ErrUserExists)
	}

	//query := `INSERT INTO users(login, email, password_hash, phone, is_active) VALUES ($1, $2, $3, $4, $5)`
	query := `INSERT INTO users(email, password_hash, is_active) VALUES ($1, $2, $3)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.PasswordHash /*user.Phone,*/, user.IsActive)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	usr, err := s.User(user.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return usr.UUID, nil
}

func (s *Storage) IsUserExists(email string) (bool, error) {
	const op = "storage.postgres.UserExists"

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRow(email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (s *Storage) User(email string) (*model.User, error) {
	const op = "storage.postgres.User"

	/*query := `
	SELECT id, login, password_hash, phone, email, is_active
	FROM users
	WHERE login=$1`
	*/
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
