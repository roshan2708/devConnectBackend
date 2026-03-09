package handlers

import (
	"devConnect/config"
	"encoding/json"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func GetFeed(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, "devconnect-session")
	currentUser := session.Values["user_id"].(string)
	query := `
	SELECT posts.id, posts.user_id, posts.content, posts.created_at
	FROM posts
	JOIN follows
	ON posts.user_id = follows.following_id
	WHERE follows.follower_id = $1
	ORDER BY posts.created_at DESC
	`
	rows, err := config.DB.Query(query, currentUser)
	if err != nil {
		http.Error(w, "Cant load the feed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Post struct {
		ID        int    `json:"id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}
	var posts []Post
	for rows.Next() {
		var post Post
		rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}
