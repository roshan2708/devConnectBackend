package handlers

import (
	"devConnect/config"
	"devConnect/middleware"
	"encoding/json"
	"net/http"
)

type LeaderboardUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AvatarURL    string `json:"avatar_url"`
	Location     string `json:"location"`
	Skills       string `json:"skills"`
	Followers    int    `json:"followers"`
	LeetcodeSolved int  `json:"leetcode"`
	CodeforcesRating int `json:"codeforces"`
	GithubStars  int    `json:"github_stars"`
	Score        int    `json:"score"`
	Streak       int    `json:"streak"`
	Posts        int    `json:"posts"`
}

// GetGlobalLeaderboard returns top users worldwide sorted by score
func GetGlobalLeaderboard(w http.ResponseWriter, r *http.Request) {
	query := `
	SELECT 
		u.id, u.name, COALESCE(u.avatar_url,''), COALESCE(u.location,''), COALESCE(u.skills,''),
		COALESCE((SELECT COUNT(*) FROM follows WHERE following_id = u.id), 0) as followers,
		COALESCE(s.leetcode_solved, 0),
		COALESCE(s.codeforces_rating, 0),
		COALESCE(s.github_commits, 0),
		COALESCE(s.total_xp, 0) as score,
		COALESCE(s.current_streak, 0),
		COALESCE((SELECT COUNT(*) FROM posts WHERE user_id = u.id), 0) as posts
	FROM users u
	LEFT JOIN user_stats s ON s.user_id = u.id
	ORDER BY score DESC
	LIMIT 50
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []LeaderboardUser
	for rows.Next() {
		var u LeaderboardUser
		rows.Scan(
			&u.ID, &u.Name, &u.AvatarURL, &u.Location, &u.Skills,
			&u.Followers, &u.LeetcodeSolved, &u.CodeforcesRating, &u.GithubStars,
			&u.Score, &u.Streak, &u.Posts,
		)
		users = append(users, u)
	}
	if users == nil {
		users = []LeaderboardUser{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetNetworkLeaderboard returns leaderboard for user's followers/following network
func GetNetworkLeaderboard(w http.ResponseWriter, r *http.Request) {
	currentUserID := middleware.GetUserID(r)
	if currentUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
	SELECT 
		u.id, u.name, COALESCE(u.avatar_url,''), COALESCE(u.location,''), COALESCE(u.skills,''),
		COALESCE((SELECT COUNT(*) FROM follows WHERE following_id = u.id), 0) as followers,
		COALESCE(s.leetcode_solved, 0),
		COALESCE(s.codeforces_rating, 0),
		COALESCE(s.github_commits, 0),
		COALESCE(s.total_xp, 0) as score,
		COALESCE(s.current_streak, 0),
		COALESCE((SELECT COUNT(*) FROM posts WHERE user_id = u.id), 0) as posts
	FROM users u
	LEFT JOIN user_stats s ON s.user_id = u.id
	WHERE u.id IN (
		SELECT following_id FROM follows WHERE follower_id = $1
		UNION
		SELECT follower_id FROM follows WHERE following_id = $1
		UNION SELECT $1
	)
	ORDER BY score DESC
	LIMIT 50
	`

	rows, err := config.DB.Query(query, currentUserID)
	if err != nil {
		http.Error(w, "Failed to fetch network leaderboard", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []LeaderboardUser
	for rows.Next() {
		var u LeaderboardUser
		rows.Scan(
			&u.ID, &u.Name, &u.AvatarURL, &u.Location, &u.Skills,
			&u.Followers, &u.LeetcodeSolved, &u.CodeforcesRating, &u.GithubStars,
			&u.Score, &u.Streak, &u.Posts,
		)
		users = append(users, u)
	}
	if users == nil {
		users = []LeaderboardUser{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
