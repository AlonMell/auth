package errors

import (
	"errors"
	"fmt"
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
	"log/slog"
	"net/http"

	serr "github.com/AlonMell/ProviderHub/internal/service/errors"
	"github.com/AlonMell/ProviderHub/pkg/logger/sl"
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
	return &Catcher{op, log, w, r}
}

func (c *Catcher) Catch(err error) {
	var errKind *serr.CustomError

	if errors.As(err, &errKind) {
		switch errKind.Kind {
		case serr.UserKind:
			//Не очень выводить ошибку пользователю которая содержить stack trace
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

	c.w.WriteHeader(http.StatusBadRequest)
	resp.WriteJSON(c.w, c.r, fmt.Errorf(InvalidRequest, err))
}
