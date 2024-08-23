package dto

import resp "providerHub/internal/lib/api/response"

type LoginRequest struct {
	Login    string `json:"login" validate:"required,alpha"`
	Password string `json:"password" validate:"required,password"`
}

type LoginResponse struct {
	resp.Response
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,alpha"`
	Password string `json:"password" validate:"required,password"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
	resp.Response
}
