package handlers

import (
	"devConnect/config"
	"devConnect/models"
	"encoding/json"
	"net/http"
)

// GetAdvancedFeed fetches the unified cross-platform feed containing posts and multi-platform events
func GetAdvancedFeed(w http.ResponseWriter, r *http.Request) {
	// Let's get the 20 most recent entries from the activity_feed table
	query := `SELECT id, user_id, activity_type, metadata, created_at FROM activity_feed ORDER BY created_at DESC LIMIT 20`
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var feed []models.ActivityFeed
	for rows.Next() {
		var item models.ActivityFeed
		if err := rows.Scan(&item.ID, &item.UserID, &item.ActivityType, &item.Metadata, &item.CreatedAt); err == nil {
			feed = append(feed, item)
		}
	}

	json.NewEncoder(w).Encode(feed)
}

// GetExtendedPosts allows filtering by advanced taxonomy (type and tags)
func GetExtendedPosts(w http.ResponseWriter, r *http.Request) {
	postType := r.URL.Query().Get("type")
	tag := r.URL.Query().Get("tag")

	query := `
		SELECT p.id, p.user_id, p.content, p.created_at, e.post_type, e.tags 
		FROM posts p
		LEFT JOIN posts_extended e ON p.id = e.post_id
		WHERE 1=1
	`
	args := []interface{}{}
	argCounter := 1

	if postType != "" {
		query += ` AND e.post_type = $` + string(rune(argCounter+'0'))
		args = append(args, postType)
		argCounter++
	}
	
	if tag != "" {
		query += ` AND e.tags LIKE $` + string(rune(argCounter+'0'))
		args = append(args, "%"+tag+"%")
		argCounter++
	}
	
	query += ` ORDER BY p.created_at DESC LIMIT 20`

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Database error finding extended posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Representing a joined view
	type AdvancedPost struct {
		ID        int    `json:"id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		Type      string `json:"type"`
		Tags      string `json:"tags"`
	}

	var results []AdvancedPost
	for rows.Next() {
		var p AdvancedPost
		var pType, pTags *string // use pointers for left join null values
		
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt, &pType, &pTags); err == nil {
			if pType != nil {
				p.Type = *pType
			} else {
				p.Type = "general"
			}
			if pTags != nil {
				p.Tags = *pTags
			}
			results = append(results, p)
		}
	}

	json.NewEncoder(w).Encode(results)
}
