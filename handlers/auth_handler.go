package handlers

import (
	"devConnect/config"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/markbates/goth/gothic"
)

func GoogleCallback(res http.ResponseWriter, req *http.Request) {

	user, err := gothic.CompleteUserAuth(res, req)

	if err != nil {
		log.Println("Auth Error:", err)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	query := `
	INSERT INTO users (id, name, email, avatar_url, last_login)
	VALUES ($1,$2,$3,$4,$5)
	ON CONFLICT (id) DO UPDATE SET
	name=EXCLUDED.name,
	avatar_url=EXCLUDED.avatar_url,
	last_login=EXCLUDED.last_login;
	`

	_, err = config.DB.Exec(
		query,
		user.UserID,
		user.Name,
		user.Email,
		user.AvatarURL,
		time.Now(),
	)

	if err != nil {
		log.Println("Insert error:", err)
	}

	session, _ := gothic.Store.Get(req, "devconnect-session")

	// Save logged-in user
	session.Values["user_id"] = user.UserID

	// Get redirect destination
	redirect, ok := session.Values["redirect"].(string)

	if !ok || redirect == "" {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:8000"
		}
		redirect = frontendURL + "/dashboard.html"
	}

	session.Save(req, res)

	http.Redirect(res, req, redirect, http.StatusTemporaryRedirect)
}
