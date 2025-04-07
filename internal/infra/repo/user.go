package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlonMell/auth/internal/domain/entity"
	vo "github.com/AlonMell/auth/internal/domain/valueObject"
	"github.com/AlonMell/auth/internal/infra/lib/logger"
	"github.com/AlonMell/grovelog/util"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AlonMell/auth/internal/domain/model"
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

func (r *UserRepo) User(
	ctx context.Context, params vo.UserParams,
) (*model.User, error) {
	ctx = logger.WithLogOp(ctx, "repo.user.User")

	builder := r.builder.
		Select("id", "email", "password_hash", "is_active").
		From("users").
		Where(sq.Eq(params))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, catch(ctx, err)
	}

	var user model.User
	err = r.db.SelectContext(ctx, &user, query, args...)
	if err != nil {
		return nil, catch(ctx, err)
	}

	return &user, nil
}

func (r *UserRepo) SaveUser(
	ctx context.Context, user model.User,
) (string, error) {
	ctx = logger.WithLogOp(ctx, "repo.user.SaveUser")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", catch(ctx, err)
	}

	txFn := func() {
		if err != nil {
			if errRb := tx.Rollback(); errRb != nil {
				err = util.WrapCtx(ctx, fmt.Errorf("error during rollback: %w", err))
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
	ctx = logger.WithLogOp(ctx, "repo.user.saveUser")

	builder := r.builder.
		Insert("users").
		Columns("id", "email", "is_active", "password_hash").
		Values(user.Id, user.Email, user.IsActive, user.PasswordHash).
		Suffix("ON CONFLICT DO NOTHING")

	query, args, err := builder.ToSql()
	if err != nil {
		return "", catch(ctx, err)
	}

	//Можно сделать с Get и Suffix("RETURNING id")
	//Можно делать с Select когда нужен полностью пользователь
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return "", catch(ctx, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return "", catch(ctx, err)
	}
	if affected == 0 {
		fmt.Println("Обосрался здесь))")
		return "", util.WrapCtx(ctx, ErrUserExists)
	}

	return user.Id, nil
}

func (r *UserRepo) DeleteUser(
	ctx context.Context, id string,
) error {
	ctx = logger.WithLogOp(ctx, "repo.user.DeleteUser")

	builder := r.builder.
		Delete("users").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return catch(ctx, err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return catch(ctx, err)
	}

	return nil
}

func (r *UserRepo) UpdateUser(
	ctx context.Context, user entity.UserMap,
) error {
	ctx = logger.WithLogOp(ctx, "repo.user.UpdateUser")

	builder := r.builder.
		Update("users").
		SetMap(user.UserParams).
		Where(sq.Eq{"id": user.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		return catch(ctx, err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return catch(ctx, err)
	}

	return nil
}
