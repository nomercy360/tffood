package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

const userIDContextKey = "user_id"

var httpErrors = map[int]string{
	http.StatusNotFound:            "Resource not found",
	http.StatusInternalServerError: "Internal server error",
	http.StatusBadRequest:          "Bad request",
	http.StatusUnauthorized:        "Unauthorized access",
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func ErrorResponder(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	if err != nil {
		log.Printf("Error: %s, Path: %s", err.Error(), r.URL.Path)
	}

	message, exists := httpErrors[statusCode]
	if !exists {
		message = "An error occurred"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				ErrorResponder(w, r, http.StatusUnauthorized, fmt.Errorf("missing auth header"))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				ErrorResponder(w, r, http.StatusUnauthorized, fmt.Errorf("invalid auth header"))
			}

			tokenString := parts[1]

			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				ErrorResponder(w, r, http.StatusUnauthorized, err)
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				ErrorResponder(w, r, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			}

			if claims.UID <= 0 {
				ErrorResponder(w, r, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
