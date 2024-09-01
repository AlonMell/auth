package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"providerHub/internal/domain/model"
)

var ErrGeneratingToken = errors.New("error generating token")

func NewToken(user *model.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.UUID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	return tokenString, nil
}
