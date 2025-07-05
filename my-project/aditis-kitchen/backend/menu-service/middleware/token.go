package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("aditi-secret"))
}

func GenerateAdminToken() (string, error) {
	return GenerateToken("aditiskitchen17@gmail.com")
}
