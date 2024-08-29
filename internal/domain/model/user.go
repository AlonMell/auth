package model

type User struct {
	UUID         string
	Email        string
	PasswordHash []byte
	//Phone        string
	IsActive bool
}
