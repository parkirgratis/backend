package config

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
    "os"
    "golang.org/x/crypto/bcrypt"
)

var JWTSecret = os.Getenv("JWT_SECRET")

func GenerateJWT(adminID string) (string, error) {
    claims := jwt.MapClaims{
        "admin_id": adminID,
        "exp":      time.Now().Add(time.Hour * 12).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(JWTSecret))
}

func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}