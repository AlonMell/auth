package user

import resp "providerHub/internal/lib/api/response"

type GetUserRequest struct {
	UUID string `json:"uuid" validate:"required,uuid"`
}

// TODO: Убрать PasswordHash из ответа

type GetUserResponse struct {
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	IsActive     bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	UUID     string `json:"uuid" validate:"uuid"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	IsActive bool   `json:"is_active"`
}

type UpdateUserResponse struct {
	resp.Response
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	IsActive bool   `json:"is_active"`
}

type CreateUserResponse struct {
	UUID string `json:"uuid"`
}

type DeleteUserRequest struct {
	UUID string `json:"uuid" validate:"required,uuid"`
}

type DeleteUserResponse struct {
	resp.Response
}
