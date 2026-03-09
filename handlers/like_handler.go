package handlers

import (
	"devConnect/config"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, "devconnect-session")
	userID := session.Values["user_id"].(string)
	vars := mux.Vars(r)
	postID := vars["postID"]
	query := `
	INSERT INTO likes (user_id, post_id)
	VALUES ($1, $2)
	ON CONFLICT (user_id, post_id) DO NOTHING
	`

	_, err := config.DB.Exec(query, userID, postID)

	if err != nil {
		http.Error(w, "Unable to like post", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Post Liked"))
}
func GetLikes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postID"]
	query := `
	SELECT user_id
	FROM likes
	WHERE post_id=$1
	`

	rows, err := config.DB.Query(query, postID)

	if err != nil {
		http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
		return

	}
	defer rows.Close()
	var likes []string

	for rows.Next() {

		var user string

		rows.Scan(&user)
		likes = append(likes, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(likes)
}
