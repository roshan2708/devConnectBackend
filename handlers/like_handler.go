package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

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

	// Find post owner
	var postOwner string

	config.DB.QueryRow(`
		SELECT user_id FROM posts WHERE id=$1
	`, postID).Scan(&postOwner)

	// Avoid notifying yourself
	if postOwner != userID {

		message := userID + " liked your post"

		config.DB.Exec(`
		INSERT INTO notifications (user_id, type, message)
		VALUES ($1,'like',$2)
		`, postOwner, message)
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

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	postID := vars["postID"]
	query := `
	DELETE FROM likes
	WHERE user_id=$1 AND post_id=$2
	`

	result, err := config.DB.Exec(query, userID, postID)

	if err != nil {
		http.Error(w, "Unable to unlike post", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Like not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("Unliked post sucessfully"))
}
