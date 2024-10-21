package middleware

import (
	"context"
	"net/http"
	"strings"
	"errors"

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
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		adminID, ok := claims["admin_id"].(string)
		if !ok {
			http.Error(w, "Invalid token claims: admin_id missing", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AdminIDContextKey, adminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
