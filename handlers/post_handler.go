package handlers

import (
	"devConnect/config"
	"encoding/json"
	"net/http"
	"time"

	"github.com/markbates/goth/gothic"
)

type CreatePostRequest struct {
	Content string `json:"content"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, "devconnect-session")
	userID := session.Values["user_id"].(string)
	var reqBody CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	query := `
	INSERT INTO posts (user_id, content, created_at)
	VALUES ($1,$2,$3)
	`
	_, err = config.DB.Exec(query, userID, reqBody.Content, time.Now())
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post Created Sucessfully"))
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	query := `
SELECT id,user_id,content,created_at
FROM posts
ORDER BY created_at DESC
`
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var userID string
		var content string
		var createdAt string
		rows.Scan(&id, &userID, &content, &createdAt)
		post := map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"content":    content,
			"created_at": createdAt,
		}
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
