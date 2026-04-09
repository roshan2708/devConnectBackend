package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {

	connStr := os.Getenv("SUPABASE_URL")

	var err error
	DB, err = sql.Open("pgx", connStr)

	if err != nil {
		log.Fatal("Unable to connect to Supabase:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}

	fmt.Println("Connected to Supabase")

	createTable()
}

func createTable() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT,
			email TEXT,
			avatar_url TEXT,
			last_login TIMESTAMP,
			bio TEXT,
			github TEXT,
			skills TEXT,
			location TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS platform_links (
			id SERIAL PRIMARY KEY,
			user_id TEXT REFERENCES users(id),
			platform TEXT,
			username TEXT,
			access_token TEXT,
			last_synced TIMESTAMP,
			UNIQUE(user_id, platform)
		);`,
		`CREATE TABLE IF NOT EXISTS developer_stats (
			user_id TEXT PRIMARY KEY REFERENCES users(id),
			total_xp INT DEFAULT 0,
			current_streak INT DEFAULT 0,
			github_commits INT DEFAULT 0,
			leetcode_solved INT DEFAULT 0,
			codeforces_rating INT DEFAULT 0,
			collaboration_score INT DEFAULT 0,
			reputation_score INT DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS developer_dna (
			user_id TEXT PRIMARY KEY REFERENCES users(id),
			primary_archetype TEXT,
			coding_style_summary TEXT,
			top_strengths TEXT,
			suggested_career_paths TEXT,
			last_analyzed TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS posts_extended (
			post_id SERIAL PRIMARY KEY,
			post_type TEXT DEFAULT 'general',
			code_snippet TEXT,
			syntax_language TEXT,
			tags TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS activity_feed (
			id SERIAL PRIMARY KEY,
			user_id TEXT REFERENCES users(id),
			activity_type TEXT,
			metadata JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS achievements (
			id SERIAL PRIMARY KEY,
			user_id TEXT REFERENCES users(id),
			badge_name TEXT,
			description TEXT,
			unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS live_sessions (
			id SERIAL PRIMARY KEY,
			host_id TEXT REFERENCES users(id),
			title TEXT,
			description TEXT,
			stream_url TEXT,
			status TEXT DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal("Could not execute query:", err, "\nQuery:", query)
		}
	}
}
