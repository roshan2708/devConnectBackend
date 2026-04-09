package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"devConnect/models"
	"encoding/json"
	"net/http"
)

// StartLiveSession creates a new Build With Me instance
func StartLiveSession(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StreamURL   string `json:"stream_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO live_sessions (host_id, title, description, stream_url) VALUES ($1, $2, $3, $4) RETURNING id`
	var sessionID int
	if err := config.DB.QueryRow(query, userID, payload.Title, payload.Description, payload.StreamURL).Scan(&sessionID); err != nil {
		http.Error(w, "Failed to start session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"session_id": sessionID, "status": "active"})
}

// GetLiveSessions returns currently active builder sessions
func GetLiveSessions(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, host_id, title, description, stream_url, status, created_at FROM live_sessions WHERE status = 'active' ORDER BY created_at DESC LIMIT 10`
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Database search failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sessions []models.LiveSession
	for rows.Next() {
		var s models.LiveSession
		if err := rows.Scan(&s.ID, &s.HostID, &s.Title, &s.Description, &s.StreamURL, &s.Status, &s.CreatedAt); err == nil {
			sessions = append(sessions, s)
		}
	}

	json.NewEncoder(w).Encode(sessions)
}
