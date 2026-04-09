package handlers

import (
	"database/sql"
	"devConnect/config"
	"devConnect/middleware"
	"devConnect/models"
	"devConnect/services/ai"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// AnalyzeDNA triggers AI processing on the user's connected accounts
func AnalyzeDNA(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	// Fetch linked GitHub username to analyze
	var githubUsername string
	err := config.DB.QueryRow(`SELECT username FROM platform_links WHERE user_id = $1 AND platform = 'github'`, userID).Scan(&githubUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "GitHub account not linked. Connect GitHub to analyze DNA.", http.StatusBadRequest)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Trigger AI Service
	dnaProfile, err := ai.AnalyzeDeveloperDNA(githubUsername)
	if err != nil {
		http.Error(w, "AI Analysis failed", http.StatusInternalServerError)
		return
	}

	// Store result
	strengthsStr := strings.Join(dnaProfile.Strengths, ",")
	careersStr := strings.Join(dnaProfile.Careers, ",")
	
	query := `
		INSERT INTO developer_dna (user_id, primary_archetype, coding_style_summary, top_strengths, suggested_career_paths, last_analyzed)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
		primary_archetype = EXCLUDED.primary_archetype,
		coding_style_summary = EXCLUDED.coding_style_summary,
		top_strengths = EXCLUDED.top_strengths,
		suggested_career_paths = EXCLUDED.suggested_career_paths,
		last_analyzed = EXCLUDED.last_analyzed
	`
	_, err = config.DB.Exec(query, userID, dnaProfile.Archetype, dnaProfile.Summary, strengthsStr, careersStr, time.Now())
	if err != nil {
		http.Error(w, "Failed to save DNA", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dnaProfile)
}

// GetDNA retrieves an existing Developer DNA
func GetDNA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	query := `SELECT primary_archetype, coding_style_summary, top_strengths, suggested_career_paths, last_analyzed 
			  FROM developer_dna WHERE user_id = $1`
			  
	row := config.DB.QueryRow(query, userID)
	
	var dna models.DeveloperDNA
	var st, car string
	dna.UserID = userID
	
	err := row.Scan(&dna.PrimaryArchetype, &dna.CodingStyleSummary, &st, &car, &dna.LastAnalyzed)
	if err != nil {
		http.Error(w, "DNA not found", http.StatusNotFound)
		return
	}
	
	dna.TopStrengths = st
	dna.SuggestedCareerPaths = car
	
	json.NewEncoder(w).Encode(dna)
}
