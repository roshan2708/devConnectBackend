package workers

import (
	"devConnect/config"
	"fmt"
	"log"
	"time"
)

// StartSyncWorker starts a background ticker that periodically fetches new data for linked platforms
func StartSyncWorker() {
	// Run every hour
	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		fmt.Println("[Sync Worker] Started background platform synchronization...")
		for range ticker.C {
			SyncPlatforms()
		}
	}()
}

func SyncPlatforms() {
	log.Println("[Sync Worker] Running platform sync iteration...")
	
	query := `SELECT user_id, platform, username FROM platform_links`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("[Sync Worker] Error querying platform links:", err)
		return
	}
	defer rows.Close()qwertyui

	for rows.Next() {
		var userID, platform, username string
		if err := rows.Scan(&userID, &platform, &username); err != nil {
			continue
		}

		if platform == "github" {
			syncGithub(userID, username)
		} else if platform == "leetcode" {
			syncLeetcode(userID, username)
		}
	}
	log.Println("[Sync Worker] Sync iteration complete.")
}

func syncGithub(userID, username string) {
	// Mock: We fetch GitHub API here and update stats
	// Example: GET https://api.github.com/users/{username}
	
	// Increment commits randomly for demonstration
	updateQuery := `
		UPDATE developer_stats 
		SET github_commits = github_commits + 5, 
		    total_xp = total_xp + 50
		WHERE user_id = $1`
		
	_, err := config.DB.Exec(updateQuery, userID)
	if err != nil {
		log.Printf("[Sync Worker] Failed updating Github stats for %s: %v\n", userID, err)
	}
}

func syncLeetcode(userID, username string) {
	// Mock: We fetch LeetCode GraphQL API here
	
	updateQuery := `
		UPDATE developer_stats 
		SET leetcode_solved = leetcode_solved + 2, 
		    total_xp = total_xp + 20
		WHERE user_id = $1`
		
	_, err := config.DB.Exec(updateQuery, userID)
	if err != nil {
		log.Printf("[Sync Worker] Failed updating Leetcode stats for %s: %v\n", userID, err)
	}
}
