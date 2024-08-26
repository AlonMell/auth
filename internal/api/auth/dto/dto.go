package dto

import resp "providerHub/internal/lib/api/response"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginResponse struct {
	resp.Response
}

type RegisterRequest struct {
	Email    string `json:"login" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
	resp.Response
}
