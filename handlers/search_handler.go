package handlers

import (
	"devConnect/config"
	"encoding/json"
	"fmt"
	"net/http"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	skill := r.URL.Query().Get("skill")
	name := r.URL.Query().Get("name")
	q := r.URL.Query().Get("q")
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

	// Support multiple query params
	if q != "" {
		skill = q
		name = q
	}

	var query string
	var args []interface{}

	if skill != "" || name != "" {
		term := skill
		if term == "" {
			term = name
		}
		query = `
		SELECT u.id, u.name, u.avatar_url, u.bio, u.skills, u.location,
		       COALESCE(u.github, '') as github,
		       COALESCE((SELECT COUNT(*) FROM follows WHERE following_id = u.id), 0) as followers,
		       COALESCE(s.leetcode_solved, 0) as leetcode_solved,
		       COALESCE(s.total_xp, 0) as total_xp
		FROM users u
		LEFT JOIN user_stats s ON s.user_id = u.id
		WHERE u.skills ILIKE $1 OR u.name ILIKE $1
		LIMIT $2 OFFSET $3
		`
		args = []interface{}{"%" + term + "%", limit, offset}
	} else {
		query = `
		SELECT u.id, u.name, u.avatar_url, u.bio, u.skills, u.location,
		       COALESCE(u.github, '') as github,
		       COALESCE((SELECT COUNT(*) FROM follows WHERE following_id = u.id), 0) as followers,
		       COALESCE(s.leetcode_solved, 0) as leetcode_solved,
		       COALESCE(s.total_xp, 0) as total_xp
		FROM users u
		LEFT JOIN user_stats s ON s.user_id = u.id
		LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserResult struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		AvatarURL    string `json:"avatar_url"`
		Bio          string `json:"bio"`
		Skills       string `json:"skills"`
		Location     string `json:"location"`
		Github       string `json:"github"`
		Followers    int    `json:"followers"`
		LeetcodeSolved int  `json:"leetcode_solved"`
		Score        int    `json:"score"`
	}

	var users []UserResult
	for rows.Next() {
		var u UserResult
		rows.Scan(&u.ID, &u.Name, &u.AvatarURL, &u.Bio, &u.Skills, &u.Location, &u.Github, &u.Followers, &u.LeetcodeSolved, &u.Score)
		users = append(users, u)
	}
	if users == nil {
		users = []UserResult{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
