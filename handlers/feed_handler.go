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
	SELECT 
posts.id,
posts.user_id,
posts.content,
posts.created_at,
COUNT(DISTINCT likes.id)*2 + COUNT(DISTINCT comments.id)*3 AS score

FROM posts

JOIN follows
ON posts.user_id = follows.following_id

LEFT JOIN likes
ON posts.id = likes.post_id

LEFT JOIN comments
ON posts.id = comments.post_id

WHERE follows.follower_id = $1

GROUP BY posts.id

ORDER BY score DESC, posts.created_at DESC

LIMIT 20
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
		Score     int    `json:"score"`
	}
	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt, &p.Score)
		posts = append(posts, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}
