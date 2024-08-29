package bcrypt

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const cost = bcrypt.DefaultCost

var (
	ErrGeneratePassword = errors.New("error generating password")
)

func GeneratePassword(password string) ([]byte, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratePassword, err)
	}
	return pass, nil
}

func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
