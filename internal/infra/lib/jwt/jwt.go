package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	AccessTTL  time.Duration `yaml:"access_ttl" env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"168h"`
	Secret     string        `yaml:"secret" env-required:"true"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var (
	ErrGeneratingToken = errors.New("error generating token")
	ErrValidatingToken = errors.New("error validation token")
)

func GenerateToken(
	id string, email string, duration time.Duration, secret string,
) (string, error) {
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	return accessToken, nil
}

func ValidateToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)

	if err != nil || !token.Valid {
		return nil, ErrValidatingToken
	}

	return claims, nil
}
