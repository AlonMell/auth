package dto

type UserGetDTO struct {
	Id string
}

type UserCreateDTO struct {
	Email    string
	Password string
	IsActive bool
}

type UserUpdateDTO struct {
	Email    string
	Password string
	IsActive bool
}

type UserDeleteDTO struct {
	Id string
}
