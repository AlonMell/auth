package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"

	"providerHub/internal/domain/model"
)

// TODO: Добавить Пул потоков

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (s *UserRepo) SaveUser(user model.User) (string, error) {
	const op = "storage.postgres.SaveUser"

	_, err := s.UserById(user.UUID)
	if err == nil {
		return "", ErrUserExists
	} else if !errors.Is(err, ErrUserNotFound) {
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

	usr, err := s.UserById(user.UUID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return usr.UUID, nil
}

func (s *UserRepo) UserById(uuid string) (*model.User, error) {
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

func (s *UserRepo) DeleteUser(uuid string) error {
	const op = "storage.postgres.DeleteUser"

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

func (s *UserRepo) UpdateUser(user model.User) error {
	const op = "storage.postgres.UpdateUser"

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

func (s *UserRepo) UserByEmail(email string) (*model.User, error) {
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
		return ErrUserNotFound
	default:
		return fmt.Errorf("%s: %w", op, err)
	}
}
