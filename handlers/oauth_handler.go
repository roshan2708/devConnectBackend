package handlers

import (
	"net/http"
	"os"

	"github.com/markbates/goth/gothic"
)

func BeginGoogleAuth(w http.ResponseWriter, r *http.Request) {

	session, _ := gothic.Store.Get(r, "devconnect-session")

	redirect := r.URL.Query().Get("redirect")

	// Default redirect
	if redirect == "" {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:8000"
		}
		redirect = frontendURL + "/dashboard.html"
	}

	session.Values["redirect"] = redirect
	session.Save(r, w)

	gothic.BeginAuthHandler(w, r)
}
