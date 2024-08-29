package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	resp "providerHub/internal/lib/api/response"
	"providerHub/internal/service"
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
	var errKind *service.CustomError

	if errors.As(err, &errKind) {
		switch errKind.Kind {
		case service.UserKind:
			resp.WriteJSON(c.w, c.r, err)
			c.w.WriteHeader(http.StatusBadRequest)
			return
		case service.InternalKind:
			c.log.Error("internal error", c.op, sl.Err(err))
			http.Error(c.w, "internal error", http.StatusInternalServerError)
			return
		case service.SystemKind:
			panic(err)
		}
	}

	resp.WriteJSON(c.w, c.r, fmt.Errorf(InvalidRequest, err))
	c.w.WriteHeader(http.StatusBadRequest)
}
