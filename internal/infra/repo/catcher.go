package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
)

func catch(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return logger.Wrap(ctx, ErrUserNotFound)
	default:
		return logger.Wrap(ctx, err)
	}
}
