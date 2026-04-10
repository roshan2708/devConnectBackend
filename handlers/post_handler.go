package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CreatePostRequest struct {
	Content string `json:"content"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
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
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := 1
	limit := 10
	if pageStr != "" {

		fmt.Sscanf(pageStr, "%d", &page)
	}
	if limitStr != "" {

		fmt.Sscanf(limitStr, "%d", limit)
	}

	offset := (page - 1) * limit
	query := `
	SELECT p.id,p.user_id,p.content,p.created_at,COALESCE(u.name, ''),COALESCE(u.avatar_url, '')
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id
	ORDER BY p.created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := config.DB.Query(query, limit, offset)
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
		var userName, avatarURL string
		rows.Scan(&id, &userID, &content, &createdAt, &userName, &avatarURL)
		post := map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"content":    content,
			"created_at": createdAt,
			"user": map[string]interface{}{
				"name":       userName,
				"avatar_url": avatarURL,
			},
		}
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	postID := vars["postID"]
	query := `
	DELETE FROM posts
	WHERE id=$1 AND user_id=$2
	`
	result, err := config.DB.Exec(query, postID, userID)
	if err != nil {
		http.Error(w, "Unable to delete the post", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Post not found or unathorised", http.StatusForbidden)
		return
	}

	w.Write([]byte("Post delted sucessfully"))

}

type UpdatePostRequest struct {
	Content string `json:"content"`
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	postID := vars["postID"]
	var body UpdatePostRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	query := `
	UPDATE posts
	SET content=$1
	WHERE id=$2 AND user_id=$3
	`

	result, err := config.DB.Exec(query, body.Content, postID, userID)

	if err != nil {
		http.Error(w, "Failed to update the psot", http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()

	if rows == 0 {
		http.Error(w, "Post not found or unauthorised", http.StatusForbidden)
		return
	}
	w.Write([]byte("Psot updated sucessfully"))
}
func GetUserPosts(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID := vars["userID"]

	query := `
	SELECT p.id, p.user_id, p.content, p.created_at, COALESCE(u.name, ''), COALESCE(u.avatar_url, '')
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id
	WHERE p.user_id=$1
	ORDER BY p.created_at DESC
	`

	rows, err := config.DB.Query(query, userID)

	if err != nil {
		http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	type User struct {
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	type Post struct {
		ID        int    `json:"id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		User      User   `json:"user"`
	}

	var posts []Post

	for rows.Next() {

		var p Post

		rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt, &p.User.Name, &p.User.AvatarURL)

		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(posts)
}

func GetTrendingPosts(w http.ResponseWriter, r *http.Request) {
	query := `
	SELECT 
	posts.id,
	posts.user_id,
	posts.content,
	posts.created_at,
	COALESCE(users.name, ''),
	COALESCE(users.avatar_url, ''),
	COUNT(DISTINCT likes.id)*2 + COUNT(DISTINCT comments.id)*3 AS score
	FROM posts
	LEFT JOIN users ON posts.user_id = users.id
	LEFT JOIN likes ON posts.id = likes.post_id
	LEFT JOIN comments ON posts.id = comments.post_id
	GROUP BY posts.id, users.name, users.avatar_url
	ORDER BY score DESC
	LIMIT 10
	`

	rows, err := config.DB.Query(query)

	if err != nil {
		http.Error(w, "Failed to get trending posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	type User struct {
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	type Post struct {
		ID        int    `json:"id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		User      User   `json:"user"`
		Score     int    `json:"score"`
	}

	var posts []Post

	for rows.Next() {

		var p Post

		rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt, &p.User.Name, &p.User.AvatarURL, &p.Score)

		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(posts)
}
