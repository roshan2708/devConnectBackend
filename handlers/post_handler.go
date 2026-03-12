package handlers

import (
	"devConnect/config"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func DeletePost(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, "devconnect-session")
	userID := session.Values["user_id"].(string)
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
	session, _ := gothic.Store.Get(r, "devconnect-session")
	userID := session.Values["user_id"].(string)

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
	session, _ := gothic.Store.Get(r, "devconncect-session")
	userID := session.Values["user_id"].(string)
	query := `
	SELECT id, user_id, content, created_at
	FROM posts
	WHERE user_id=$1
	ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to get psots", http.StatusInternalServerError)
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
		var p Post
		rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt)
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
	COUNT(DISTINCT likes.id)*2 + COUNT(DISTINCT comments.id)*3 AS score
	FROM posts
	LEFT JOIN likes ON posts.id = likes.post_id
	LEFT JOIN comments ON posts.id = comments.post_id
	GROUP BY posts.id
	ORDER BY score DESC
	LIMIT 10
	`

	rows, err := config.DB.Query(query)

	if err != nil {
		http.Error(w, "Failed to get trending posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Post struct {
		ID        int    `json:"id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		Score     int    `json:"score"`
	}

	var posts []Post

	for rows.Next() {

		var p Post

		rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt, &p.Score)

		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(posts)
}
