package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

func GenerateJWT(userId string) (string, error) {
	JwtSecret := os.Getenv("JWT_SECRET")

	clams := &Claims{
		jwt.RegisteredClaims{
			Issuer:    userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clams)

	return token.SignedString([]byte(JwtSecret))
}
