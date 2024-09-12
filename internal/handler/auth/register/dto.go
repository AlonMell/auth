package register

import (
	resp "github.com/AlonMell/ProviderHub/internal/infra/lib/api/response"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Response struct {
	Id string `json:"id"`
	resp.Response
}
