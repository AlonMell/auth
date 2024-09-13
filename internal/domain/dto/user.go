package dto

type UserGetReq struct {
	Id string `json:"id" validate:"required,uuid"`
}

type UserCreateReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	IsActive bool   `json:"is_active"`
}

type UserUpdateReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	IsActive bool   `json:"is_active"`
}

type UserDeleteReq struct {
	Id string `json:"id" validate:"required,uuid"`
}
