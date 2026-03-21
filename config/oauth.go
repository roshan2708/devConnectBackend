package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config

func InitOauth() {
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:3000/auth/google/callback"
	}

	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
