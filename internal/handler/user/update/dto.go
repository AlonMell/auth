package update

import (
	resp "providerHub/internal/infra/lib/api/response"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	IsActive bool   `json:"is_active"`
}

type Response struct {
	resp.Response
}
