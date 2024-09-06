package auth

import resp "providerHub/internal/lib/api/response"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type RegisterResponse struct {
	UUID string `json:"uuid"`
	resp.Response
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
