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

	redirectURL := fmt.Sprintf("devconnect://auth?userId=%s&name=%s&success=true", user.UserID, user.Name)

	res.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(res, `
		<html>
			<body style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; font-family: sans-serif; text-align: center; padding: 20px;">
				<h2 style="color: #4285F4;">Login Successful!</h2>
				<p>You are being redirected back to the DevConnect app.</p>
				<p style="font-size: 0.9em; color: #666;">If you are not redirected automatically, click the button below:</p>
				<a href="%s" style="display: inline-block; padding: 12px 24px; background: #4285F4; color: white; text-decoration: none; border-radius: 8px; margin-top: 20px; font-weight: bold; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">Return to App</a>
				<script>
					// Attempt automatic redirect
					window.location.href = "%s";
					
					// Auto-close or redirect fallback after a delay
					setTimeout(function() {
						window.location.href = "%s";
					}, 2000);
				</script>
			</body>
		</html>
	`, redirectURL, redirectURL, redirectURL)
}
