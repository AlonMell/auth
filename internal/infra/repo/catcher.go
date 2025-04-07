package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlonMell/grovelog/util"
)

func catch(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return util.WrapCtx(ctx, ErrUserNotFound)
	default:
		return util.WrapCtx(ctx, err)
	}
}
