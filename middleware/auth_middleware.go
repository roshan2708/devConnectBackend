package middleware

import (
	"net/http"

	"github.com/markbates/goth/gothic"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for X-User-ID header first (for mobile apps)
		userID := r.Header.Get("X-User-ID")
		if userID != "" {
			next.ServeHTTP(w, r)
			return
		}

		session, err := gothic.Store.Get(r, "devconnect-session")
		if err != nil {
			http.Error(w, "Session-error", http.StatusInternalServerError)
			return
		}
		sessionUserID, ok := session.Values["user_id"].(string)
		if !ok || sessionUserID == "" {
			http.Error(w, "Unauthorised", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	})
}
