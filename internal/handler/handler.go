package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	resp "providerHub/internal/lib/api/response"
	serr "providerHub/internal/service/errors"
	"providerHub/pkg/logger/sl"
)

const (
	InvalidRequest = "invalid request %w"
)

type Catcher struct {
	op  string
	log *slog.Logger
	w   http.ResponseWriter
	r   *http.Request
}

func NewCatcher(op string, log *slog.Logger, w http.ResponseWriter, r *http.Request) *Catcher {
	return &Catcher{
		op:  op,
		log: log,
		w:   w,
		r:   r,
	}
}

func (c *Catcher) Catch(err error) {
	var errKind *serr.CustomError

	if errors.As(err, &errKind) {
		switch errKind.Kind {
		case serr.UserKind:
			resp.WriteJSON(c.w, c.r, err)
			c.w.WriteHeader(errKind.Code)
			return
		case serr.InternalKind:
			c.log.Error("internal error", c.op, sl.Err(err))
			http.Error(c.w, "internal error", errKind.Code)
			return
		case serr.SystemKind:
			panic(err)
		}
	}

	resp.WriteJSON(c.w, c.r, fmt.Errorf(InvalidRequest, err))
	c.w.WriteHeader(http.StatusBadRequest)
}
