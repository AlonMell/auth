package get

type Request struct {
	Id string `json:"id" validate:"required,uuid"`
}

// TODO: Убрать PasswordHash из ответа

type Response struct {
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	IsActive     bool   `json:"is_active"`
}
