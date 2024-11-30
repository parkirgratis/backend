package middleware

import (
    "context"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/gocroot/config"
    "go.mongodb.org/mongo-driver/bson"
)

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

        claims, err := validateJWT(tokenString)
        if err != nil {
            if err == jwt.ErrTokenExpired {
                fmt.Println("Token is expired")
                http.Error(w, "Token expired", http.StatusUnauthorized)

                go deleteExpiredToken(tokenString)
                return
            }

            fmt.Println("Invalid token:", err)
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        adminID, ok := claims["admin_id"].(string)
        if !ok {
            fmt.Println("admin_id not found in token claims")
            http.Error(w, "Invalid token claims: admin_id missing", http.StatusUnauthorized)
            return
        }
        fmt.Println("admin_id from token:", adminID)

        r.Header.Set("admin_id", adminID)


        next.ServeHTTP(w, r)
    })
}

func validateJWT(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.JWTSecret), nil
    })

    if err != nil {
        var ve *jwt.ValidationError
        if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
            return nil, jwt.ErrTokenExpired
        }
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, fmt.Errorf("invalid token claims")
    }

    return claims, nil
}

func deleteExpiredToken(tokenString string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"token": tokenString}
    _, err := config.Mongoconn.Collection("tokens").DeleteOne(ctx, filter)
    if err != nil {
        fmt.Printf("Failed to delete expired token: %v\n", err)
    } else {
        fmt.Println("Expired token deleted successfully")
    }
}
