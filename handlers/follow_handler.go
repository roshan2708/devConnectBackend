package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetUserID(r)

	vars := mux.Vars(r)
	targetUser := vars["userID"]

	query := `
	INSERT INTO follows (follower_id, following_id)
	VALUES ($1,$2)
	`

	_, err := config.DB.Exec(query, currentUser, targetUser)

	if err != nil {
		http.Error(w, "Unable to follow", http.StatusInternalServerError)
		return
	}

	// Create notification
	notifQuery := `
	INSERT INTO notifications (user_id, type, message)
	VALUES ($1, $2, $3)
	`

	message := currentUser + " started following you"

	_, err = config.DB.Exec(notifQuery, targetUser, "follow", message)

	if err != nil {
		// We don't break the request if notification fails
		// Just log it
		fmt.Println("Notification error:", err)
	}

	w.Write([]byte("Followed Successfully"))
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	query := `
	SELECT u.id, COALESCE(u.name,''), COALESCE(u.avatar_url,''), COALESCE(u.location,''), COALESCE(u.bio,'')
	FROM follows f
	JOIN users u ON u.id = f.follower_id
	WHERE f.following_id=$1
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Unable to get followers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserInfo struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Location  string `json:"location"`
		Bio       string `json:"bio"`
	}
	var followers []UserInfo
	for rows.Next() {
		var u UserInfo
		rows.Scan(&u.ID, &u.Name, &u.AvatarURL, &u.Location, &u.Bio)
		followers = append(followers, u)
	}
	if followers == nil {
		followers = []UserInfo{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	query := `
	SELECT u.id, COALESCE(u.name,''), COALESCE(u.avatar_url,''), COALESCE(u.location,''), COALESCE(u.bio,'')
	FROM follows f
	JOIN users u ON u.id = f.following_id
	WHERE f.follower_id=$1
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch following", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserInfo struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Location  string `json:"location"`
		Bio       string `json:"bio"`
	}
	var following []UserInfo
	for rows.Next() {
		var u UserInfo
		rows.Scan(&u.ID, &u.Name, &u.AvatarURL, &u.Location, &u.Bio)
		following = append(following, u)
	}
	if following == nil {
		following = []UserInfo{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetUserID(r)
	vars := mux.Vars(r)

	targetUser := vars["userID"]
	query := `
	DELETE FROM follows
	WHERE follower_id=$1 AND following_id=$2
	`

	result, err := config.DB.Exec(query, currentUser, targetUser)

	if err != nil {
		http.Error(w, "Unable to unfloolw user", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Follow relationship not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Unfollowed user sucessully"))
}
