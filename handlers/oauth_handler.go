package handlers

import (
	"net/http"

	"github.com/markbates/goth/gothic"
)

func BeginGoogleAuth(w http.ResponseWriter, r *http.Request) {

	session, _ := gothic.Store.Get(r, "devconnect-session")

	redirect := r.URL.Query().Get("redirect")

	if redirect != "" {
		session.Values["redirect"] = redirect
		session.Save(r, w)
	}

	gothic.BeginAuthHandler(w, r)
}
