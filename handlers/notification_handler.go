package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}

	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	offset := (page - 1) * limit

	query := `
	SELECT id,type,message,created_at,read
	FROM notifications
	WHERE user_id=$1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := config.DB.Query(query, userID, limit, offset)
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
func MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	notificationID := vars["id"]
	query := `
	UPDATE notifications
	SET read = true
	WHERE id=$1 AND user_id=$2
	`

	result, err := config.DB.Exec(query, notificationID, userID)

	if err != nil {
		http.Error(w, "Unable to read notifications", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()

	if rows == 0 {
		http.Error(w, "Notifiction not found", http.StatusForbidden)
		return
	}

	w.Write([]byte("Marked notification as read"))
}
