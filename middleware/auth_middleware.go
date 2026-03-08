package middleware

import (
	"net/http"

	"github.com/markbates/goth/gothic"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := gothic.Store.Get(r, "devconnect-session")
		if err != nil {
			http.Error(w, "Session-error", http.StatusInternalServerError)
			return
		}
		userID, ok := session.Values["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorised", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	})
}
