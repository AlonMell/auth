package catcher

import (
	"context"
	"errors"
	"fmt"
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	serr "github.com/AlonMell/ProviderHub/internal/service/catcher"
	"log/slog"
	"net/http"

	"github.com/AlonMell/ProviderHub/pkg/logger/sl"
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
			c.log.ErrorContext(logger.ErrorCtx(c.ctx, err), "user kind error", sl.Err(err))
			resp.WriteJSON(c.w, c.r, err.Error())
			return
		case serr.InternalKind:
			c.log.ErrorContext(logger.ErrorCtx(c.ctx, err), "internal kind error", sl.Err(err))
			resp.WriteJSON(c.w, c.r, "internal error")
			return
		case serr.SystemKind:
			c.log.ErrorContext(logger.ErrorCtx(c.ctx, err), "system kind error", sl.Err(err))
			panic(err)
		}
	}

	resp.Status(c.r, http.StatusBadRequest)
	resp.WriteJSON(c.w, c.r, fmt.Sprintf(InvalidRequest, err.Error()))
}
