package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ProfileUpdate struct {
	Bio      string `json:"bio"`
	Github   string `json:"github"`
	Skills   string `json:"skills"`
	Location string `json:"location"`
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var profile ProfileUpdate
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	query := `
	UPDATE users
	SET bio=$1, github=$2, skills=$3, location=$4
	WHERE id=$5
	`

	_, err = config.DB.Exec(
		query,
		profile.Bio,
		profile.Github,
		profile.Skills,
		profile.Location,
		userID,
	)

	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Profile updated"))

}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	query := `
	SELECT id,name,email,avatar_url,bio,github,skills,location
	FROM users
	WHERE id=$1
	`
	row := config.DB.QueryRow(query, userID)

	type User struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar_url"`
		Bio      string `json:"bio"`
		Github   string `json:"github"`
		Skills   string `json:"skills"`
		Location string `json:"location"`
	}

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Avatar,
		&user.Bio,
		&user.Github,
		&user.Skills,
		&user.Location,
	)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)

}
