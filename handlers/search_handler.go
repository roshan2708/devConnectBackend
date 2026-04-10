package handlers

import (
	"devConnect/config"
	"encoding/json"
	"fmt"
	"net/http"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	skill := r.URL.Query().Get("skill")
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
	SELECT id, name, skills
	FROM users
	WHERE skills ILIKE $1
	LIMIT $2 OFFSET $3
	`

	rows, err := config.DB.Query(query, "%"+skill+"%", limit, offset)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type User struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Skills string `json:"skills"`
	}

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Skills)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}
