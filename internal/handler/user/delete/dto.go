package delete

import (
	resp "providerHub/internal/infra/lib/api/response"
)

type Request struct {
	Id string `json:"id" validate:"required,uuid"`
}

type Response struct {
	resp.Response
}
