package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gocroot/config"
)

type contextKey string

const AdminIDContextKey contextKey = "admin_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			fmt.Println("Authorization header missing")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			fmt.Println("Invalid token format")
			return
		}

		tokenString := parts[1]
		fmt.Println("Token received:", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					fmt.Println("Token is malformed")
					http.Error(w, "Malformed token", http.StatusUnauthorized)
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					fmt.Println("Token is expired")
					http.Error(w, "Token expired", http.StatusUnauthorized)
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					fmt.Println("Token not valid yet")
					http.Error(w, "Token not valid yet", http.StatusUnauthorized)
				} else {
					fmt.Println("Invalid token:", err)
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				}
			} else {
				fmt.Println("Error parsing token:", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			return
		}

		if !token.Valid {
			fmt.Println("Token is not valid")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			fmt.Println("Invalid token claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		adminID, ok := claims["admin_id"].(string)
		if !ok {
			fmt.Println("admin_id not found in token claims")
			http.Error(w, "Invalid token claims: admin_id missing", http.StatusUnauthorized)
			return
		}
		fmt.Println("admin_id from token:", adminID)

		ctx := context.WithValue(r.Context(), AdminIDContextKey, adminID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
