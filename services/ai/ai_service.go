package ai

import (
	"fmt"
	"log"
)

// DNAProfile represents the output from the AI
type DNAProfile struct {
	Archetype string   `json:"primary_archetype"`
	Summary   string   `json:"coding_style_summary"`
	Strengths []string `json:"top_strengths"`
	Careers   []string `json:"suggested_career_paths"`
}

// AnalyzeDeveloperDNA mocks an integration with Gemini/OpenAI
// In production, this would pass github commit data to the LLM and return a structured JSON response
func AnalyzeDeveloperDNA(githubUsername string) (*DNAProfile, error) {
	log.Printf("Analyzing Developer DNA for GitHub user: %s (AI Mock)\n", githubUsername)

	// Mock decision heuristic based on simple variables
	archetype := "Problem Solver"
	if len(githubUsername) > 6 {
		archetype = "Builder"
	}

	dna := &DNAProfile{
		Archetype: archetype,
		Summary:   fmt.Sprintf("User %s exhibits consistent backend-focused architectural abilities.", githubUsername),
		Strengths: []string{"System Design", "Go", "API Development"},
		Careers:   []string{"Senior Backend Engineer", "Platform Architect"},
	}

	return dna, nil
}

func AnalyzeOpportunityReadiness(commits int, posts int, likes int) int {
	// A simple heuristic determining how "collaboration ready" someone is.
	// Max score 100.
	score := commits/10 + posts*2 + likes
	if score > 100 {
		score = 100
	}
	return score
}
