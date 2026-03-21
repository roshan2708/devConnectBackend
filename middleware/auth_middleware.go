package middleware

import (
	"context"
	"devConnect/utils"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Check for X-User-ID header (legacy/mobile support)
		userID := r.Header.Get("X-User-ID")
		if userID != "" {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 2. Check for Bearer Token in Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				validUserID, err := utils.ValidateToken(tokenString)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserIDKey, validUserID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		http.Error(w, "Unauthorised", http.StatusUnauthorized)
	})
}

func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
