package handlers

import (
	"devConnect/config"
	"encoding/json"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, "devconnect-session")
	userID := session.Values["user_id"].(string)
	query := `
	SELECT id,type,message,created_at,read
	FROM notifications
	WHERE user_id=$1
	ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Unable to laod notifications", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	type Notification struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
		Read      bool   `json:"read"`
	}
	var notifications []Notification
	for rows.Next() {
		var n Notification
		rows.Scan(&n.ID, &n.Type, &n.Message, &n.CreatedAt, &n.Read)
		notifications = append(notifications, n)

	}
	json.NewEncoder(w).Encode(notifications)
}
