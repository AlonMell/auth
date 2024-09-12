package delete

import (
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
)

type Request struct {
	Id string `json:"id" validate:"required,uuid"`
}

type Response struct {
	resp.Response
}
