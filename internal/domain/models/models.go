package models

type User struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	IsActive     bool   `json:"is_active"`
}
