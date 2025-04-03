package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/utils"
)

type contextKey string

const (
	ContextUserID   contextKey = "userID"
	ContextUserRole contextKey = "userRole"
)

func JWTAuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
			claims, err := utils.ValidateJWT(tokenStr, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextUserRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
