package register

import resp "providerHub/internal/lib/api/response"

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Response struct {
	Id string `json:"id"`
	resp.Response
}
