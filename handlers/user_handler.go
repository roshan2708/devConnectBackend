package handlers

import (
	"devConnect/config"
	"encoding/json"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Check for X-User-ID header first (for mobile apps)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		session, err := gothic.Store.Get(r, "devconnect-session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		var ok bool
		userID, ok = session.Values["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "User not logged in ", http.StatusUnauthorized)
			return
		}
	}
	query := `
	SELECT id,name,email,avatar_url
	FROM users
	WHERE id=$1
	`
	var user struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	err := config.DB.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.AvatarURL)
	if err != nil {
		http.Error(w, "User not found ", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}
