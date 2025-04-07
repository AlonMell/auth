package errors

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	resp "github.com/AlonMell/auth/internal/infra/lib/api/response"
	serr "github.com/AlonMell/auth/internal/service/errors"
	"github.com/AlonMell/grovelog/util"
)

const (
	InvalidRequest = "invalid request: %s"
)

type Catcher struct {
	ctx context.Context
	log *slog.Logger
	w   http.ResponseWriter
	r   *http.Request
}

func NewCatcher(
	ctx context.Context,
	log *slog.Logger,
	w http.ResponseWriter,
	r *http.Request,
) *Catcher {
	return &Catcher{ctx, log, w, r}
}

func (c *Catcher) Catch(err error) {
	var errKind *serr.ServiceError

	if errors.As(err, &errKind) {
		resp.Status(c.r, errKind.Code)

		switch errKind.Kind {
		case serr.UserKind:
			resp.WriteJSON(c.w, c.r, err.Error())
			return
		case serr.InternalKind:
			c.log.ErrorContext(util.ErrorCtx(c.ctx, err), "internal error", util.Err(err))
			resp.WriteJSON(c.w, c.r, "internal error")
			return
		case serr.SystemKind:
			panic(err)
		}
	}

	resp.Status(c.r, http.StatusBadRequest)
	resp.WriteJSON(c.w, c.r, fmt.Sprintf(InvalidRequest, err.Error()))
}
