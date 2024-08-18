package model

type User struct {
	ID           string
	Login        string
	PasswordHash []byte
	//Phone        string
	//Email        string
	IsActive bool
}
