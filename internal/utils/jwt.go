//go:generate mockgen -source ./jwt.go -destination=./mocks/jwt.go -package=mock_jwt
package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTGenerator interface {
	GenerateJWT(role string) (string, error)
}

type JWTGen struct{}

func (j *JWTGen) GenerateJWT(role string) (string, error) {
	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		return "", errors.New("environment variable JWT_SECRET not found")
	}
	jwtSecret := []byte(jwt)

	return token.SignedString(jwtSecret)
}
