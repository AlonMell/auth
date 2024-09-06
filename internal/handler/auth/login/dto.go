package login

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Response struct {
	Jwt string `json:"jwt"`
}
