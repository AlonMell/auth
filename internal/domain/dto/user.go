package dto

type UserGet struct {
	Id string
}

type UserCreate struct {
	Email    string
	Password string
	IsActive bool
}

type UserUpdate struct {
	Email    string
	Password string
	IsActive bool
}

type UserDelete struct {
	Id string
}
