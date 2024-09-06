package post

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	IsActive bool   `json:"is_active"`
}

type Response struct {
	Id string `json:"id"`
}
