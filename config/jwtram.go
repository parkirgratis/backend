package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JWTSecrets = os.Getenv("JWT_SECRET")

func GenerateJWTs(adminID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": adminID,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}
