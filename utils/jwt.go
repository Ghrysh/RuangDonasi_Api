package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, role string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("fallback_secret_key")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}