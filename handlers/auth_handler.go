package handlers

import (
	"devConnect/config"
	"fmt"
	"log"
	"net/http"
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

	_, err = config.DB.Exec(query,
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
	session.Values["user_id"] = user.UserID
	session.Save(req, res)

	redirectURL := fmt.Sprintf("/?success=true&name=%s", user.Name)

	http.Redirect(res, req, redirectURL, http.StatusTemporaryRedirect)
}
