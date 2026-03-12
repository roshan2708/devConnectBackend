package handlers

import (
	"devConnect/config"
	"encoding/json"
	"net/http"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	skill := r.URL.Query().Get("skill")
	query := `
	SELECT id, name, skills
	FROM users
	WHERE skills ILIKE $1
	`

	rows, err := config.DB.Query(query, "%"+skill+"%")
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
