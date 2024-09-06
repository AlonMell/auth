package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"

	"providerHub/internal/domain/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

const (
	userByIdQuery    = `SELECT id, email, password_hash, is_active FROM users WHERE id=$1`
	userByEmailQuery = `SELECT id, email, password_hash, is_active FROM users WHERE email=$1`
	deleteUserQuery  = `DELETE FROM users WHERE id=$1`
	updateUserQuery  = `UPDATE users SET email=$1, password_hash=$2, is_active=$3 WHERE id=$4`
	insertUserQuery  = `INSERT INTO users(id, email, password_hash, is_active) VALUES ($1, $2, $3, $4)`
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (s *UserRepo) SaveUser(ctx context.Context, user model.User) (string, error) {
	const op = "storage.postgres.SaveUser"

	_, err := s.UserById(ctx, user.Id)
	if err == nil {
		return "", ErrUserExists
	} else if !errors.Is(err, ErrUserNotFound) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare(insertUserQuery)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Email, user.PasswordHash /*user.Phone,*/, user.IsActive)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	usr, err := s.UserById(ctx, user.Id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return usr.Id, nil
}

func (s *UserRepo) UserById(
	ctx context.Context, id string,
) (*model.User, error) {
	const op = "storage.postgres.UserById"

	stmt, err := s.db.Prepare(userByIdQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.IsActive)
	if err != nil {
		return nil, errorHandler(err, op)
	}

	return &user, nil
}

func (s *UserRepo) DeleteUser(
	ctx context.Context, id string,
) error {
	const op = "storage.postgres.DeleteUser"

	stmt, err := s.db.Prepare(deleteUserQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *UserRepo) UpdateUser(
	ctx context.Context, user model.User,
) error {
	const op = "storage.postgres.UpdateUser"

	stmt, err := s.db.Prepare(updateUserQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.PasswordHash, user.IsActive, user.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *UserRepo) UserByEmail(
	ctx context.Context, email string,
) (*model.User, error) {
	const op = "storage.postgres.User"

	stmt, err := s.db.Prepare(userByEmailQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.IsActive)
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
