package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"providerHub/internal/domain/model"
	"time"
)

// TODO: Разобраться как генерируется jwt

func NewToken(user model.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	//claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
