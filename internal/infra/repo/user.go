package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AlonMell/ProviderHub/internal/domain/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepo struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func NewUserRepo(db *sqlx.DB, placeHolder sq.PlaceholderFormat) *UserRepo {
	return &UserRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(placeHolder),
	}
}

func (r *UserRepo) SaveUser(
	ctx context.Context, user model.User,
) (string, error) {
	const op = "repo.user.SaveUser"
	ctx = logger.WithLogOp(ctx, op)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", errorHandler(op, err)
	}

	txFn := func() {
		if err != nil {
			if errRb := tx.Rollback(); errRb != nil {
				err = fmt.Errorf("error during rollback: %w", err)
			}
			return
		}
		err = tx.Commit()
	}

	defer txFn()

	return r.saveUser(tx, ctx, user)
}

func (r *UserRepo) saveUser(
	tx *sqlx.Tx, ctx context.Context, user model.User,
) (string, error) {
	const op = "repo.user.saveUser"
	ctx = logger.WithLogOp(ctx, op)

	builder := r.builder.
		Insert("users").
		Values(user.Id, user.Email, user.IsActive, user.PasswordHash).
		Suffix("ON CONFLICT DO NOTHING")

	query, args, err := builder.ToSql()
	if err != nil {
		return "", errorHandler(op, err)
	}

	//Можно сделать с Get и Suffix("RETURNING id")
	//Можно делать с Select когда нужен полностью пользователь
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return "", errorHandler(op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return "", errorHandler(op, err)
	}
	if affected == 0 {
		return "", ErrUserExists
	}

	return user.Id, nil
}

func (r *UserRepo) UserById(
	ctx context.Context, id string,
) (*model.User, error) {
	const op = "repo.user.UserById"
	ctx = logger.WithLogOp(ctx, op)

	builder := r.builder.
		Select("id", "email", "password_hash", "is_active").
		From("users").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errorHandler(op, err)
	}

	var user model.User
	err = r.db.SelectContext(ctx, &user, query, args...)
	if err != nil {
		return nil, errorHandler(op, err)
	}

	return &user, nil
}

func (r *UserRepo) DeleteUser(
	ctx context.Context, id string,
) error {
	const op = "repo.user.DeleteUser"
	ctx = logger.WithLogOp(ctx, op)

	builder := r.builder.
		Delete("users").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return errorHandler(op, err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errorHandler(op, err)
	}

	return nil
}

func (r *UserRepo) UpdateUser(
	ctx context.Context, user model.User,
) error {
	const op = "repo.user.UpdateUser"
	ctx = logger.WithLogOp(ctx, op)

	builder := r.builder.
		Update("users").
		Set("email", user.Email).
		Set("password_hash", user.PasswordHash).
		Set("is_active", user.IsActive).
		Where(sq.Eq{"id": user.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		return errorHandler(op, err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errorHandler(op, err)
	}

	return nil
}

func (r *UserRepo) UserByEmail(
	ctx context.Context, email string,
) (*model.User, error) {
	const op = "repo.user.UserByEmail"
	ctx = logger.WithLogOp(ctx, op)

	builder := r.builder.
		Select("id", "email", "password_hash", "is_active").
		From("users").
		Where(sq.Eq{"email": email})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errorHandler(op, err)
	}

	var user model.User
	err = r.db.SelectContext(ctx, &user, query, args...)
	if err != nil {
		return nil, errorHandler(op, err)
	}

	ctx = logger.WithLogUserID(ctx, user.Id)

	return &user, nil
}

func errorHandler(op string, err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrUserNotFound
	default:
		return fmt.Errorf("%s: %w", op, err)
	}
}
