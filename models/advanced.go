package models

import "time"

// PlatformLink represents linked third-party accounts
type PlatformLink struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	Platform    string    `json:"platform"`
	Username    string    `json:"username"`
	AccessToken string    `json:"access_token,omitempty"`
	LastSynced  time.Time `json:"last_synced"`
}

// DeveloperStats holds aggregated gamification and stat values
type DeveloperStats struct {
	UserID             string `json:"user_id"`
	TotalXP            int    `json:"total_xp"`
	CurrentStreak      int    `json:"current_streak"`
	GithubCommits      int    `json:"github_commits"`
	LeetcodeSolved     int    `json:"leetcode_solved"`
	CodeforcesRating   int    `json:"codeforces_rating"`
	CollaborationScore int    `json:"collaboration_score"`
	ReputationScore    int    `json:"reputation_score"`
}

// DeveloperDNA stores AI generated insights
type DeveloperDNA struct {
	UserID               string    `json:"user_id"`
	PrimaryArchetype     string    `json:"primary_archetype"`
	CodingStyleSummary   string    `json:"coding_style_summary"`
	TopStrengths         string    `json:"top_strengths"` // Typically parsed as JSON list
	SuggestedCareerPaths string    `json:"suggested_career_paths"` // Typically parsed as JSON list
	LastAnalyzed         time.Time `json:"last_analyzed"`
}

// PostExtended adds extra context onto standard Posts
type PostExtended struct {
	PostID         int    `json:"post_id"`
	PostType       string `json:"post_type"`
	CodeSnippet    string `json:"code_snippet"`
	SyntaxLanguage string `json:"syntax_language"`
	Tags           string `json:"tags"` // JSON array string
}

// ActivityFeed represents cross-platform activity timeline
type ActivityFeed struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	ActivityType string    `json:"activity_type"`
	Metadata     string    `json:"metadata"` // JSONB from DB
	CreatedAt    time.Time `json:"created_at"`
}

// Achievement tracks badges and milestones
type Achievement struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	BadgeName   string    `json:"badge_name"`
	Description string    `json:"description"`
	UnlockedAt  time.Time `json:"unlocked_at"`
}

// LiveSession defines "Build With Me" sessions
type LiveSession struct {
	ID          int       `json:"id"`
	HostID      string    `json:"host_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StreamURL   string    `json:"stream_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
