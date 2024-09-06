package model

import "github.com/google/uuid"

type User struct {
	Id           string
	Email        string
	PasswordHash []byte
	//Phone        string
	IsActive bool
}

func NewUser(email string, passwordHash []byte, isActive bool) *User {
	return &User{
		Id:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		IsActive:     isActive,
	}
}
