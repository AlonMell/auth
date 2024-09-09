package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

const cost = bcrypt.DefaultCost

func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
