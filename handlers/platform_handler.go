package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"devConnect/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// LinkPlatform links a third party platform (GitHub/Leetcode) to a user account
func LinkPlatform(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	platform := vars["platform"] // 'github' or 'leetcode'

	var payload map[string]string
	json.NewDecoder(r.Body).Decode(&payload)

	username := payload["username"]
	token := payload["access_token"] // Optional 

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO platform_links (user_id, platform, username, access_token, last_synced)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, platform) DO UPDATE 
		SET username = EXCLUDED.username, access_token = EXCLUDED.access_token, last_synced = EXCLUDED.last_synced
	`
	_, err := config.DB.Exec(query, userID, platform, username, token, time.Now())
	if err != nil {
		http.Error(w, "Failed to link platform", http.StatusInternalServerError)
		return
	}

	// Make sure they have a row in developer_stats
	config.DB.Exec(`INSERT INTO developer_stats (user_id) VALUES ($1) ON CONFLICT DO NOTHING`, userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully linked " + platform})
}

// GetUserStats retrieves the advanced gamification data for a user
func GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	query := `SELECT total_xp, current_streak, github_commits, leetcode_solved, codeforces_rating, collaboration_score, reputation_score 
			  FROM developer_stats WHERE user_id = $1`
			  
	row := config.DB.QueryRow(query, userID)
	
	var stats models.DeveloperStats
	stats.UserID = userID
	err := row.Scan(&stats.TotalXP, &stats.CurrentStreak, &stats.GithubCommits, &stats.LeetcodeSolved, &stats.CodeforcesRating, &stats.CollaborationScore, &stats.ReputationScore)
	
	if err != nil {
		// If no stats yet, return zeros (instead of 404)
		stats = models.DeveloperStats{UserID: userID}
	}
	
	json.NewEncoder(w).Encode(stats)
}
