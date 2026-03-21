package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CommentRequest struct {
	Content string `json:"content"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	postID := vars["postID"]

	var body CommentRequest
	json.NewDecoder(r.Body).Decode(&body)

	query := `
	INSERT INTO comments (post_id, user_id, content)
	VALUES ($1,$2,$3)
	`

	_, err := config.DB.Exec(query, postID, userID, body.Content)

	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	// Find post owner
	var postOwner string

	config.DB.QueryRow(`
		SELECT user_id FROM posts WHERE id=$1
	`, postID).Scan(&postOwner)

	if postOwner != userID {

		message := userID + " commented on your post"

		config.DB.Exec(`
		INSERT INTO notifications (user_id,type,message)
		VALUES ($1,'comment',$2)
		`, postOwner, message)
	}

	w.Write([]byte("Added comment successfully"))
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postID"]

	query := `
	SELECT user_id, content, created_at
	FROM comments
	WHERE post_id=$1
	ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(query, postID)

	if err != nil {
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	type Comment struct {
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}
	var comments []Comment

	for rows.Next() {

		var c Comment
		rows.Scan(&c.UserID, &c.Content, &c.CreatedAt)
		comments = append(comments, c)
	}
	json.NewEncoder(w).Encode(comments)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)

	commentID := vars["commentID"]
	query := `
	DELETE FROM comments
	WHERE id=$1 AND user_id=$2
	`

	result, err := config.DB.Exec(query, commentID, userID)

	if err != nil {
		http.Error(w, "Unable to delete comment", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Comment not found or unauthorised", http.StatusForbidden)
		return
	}

	w.Write([]byte("Comment deleted sucessfully"))
}
