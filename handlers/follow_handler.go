package handlers

import (
	"devConnect/config"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {

	session, _ := gothic.Store.Get(r, "devconnect-session")
	currentUser := session.Values["user_id"].(string)

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
	userID := vars["userID"] // FIXED

	query := `
	SELECT follower_id
	FROM follows
	WHERE following_id=$1
	`

	rows, err := config.DB.Query(query, userID)

	if err != nil {
		http.Error(w, "Unable to get followers", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var followers []string

	for rows.Next() {

		var follower string

		rows.Scan(&follower)

		followers = append(followers, follower)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID := vars["userID"] // FIXED

	query := `
	SELECT following_id
	FROM follows
	WHERE follower_id=$1
	`

	rows, err := config.DB.Query(query, userID)

	if err != nil {
		http.Error(w, "Failed to fetch following", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var following []string

	for rows.Next() {

		var user string

		rows.Scan(&user)

		following = append(following, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, "devocnnect-session")
	currentUser := session.Values["user_id"].(string)
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
