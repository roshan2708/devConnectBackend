package handlers

import (
	"context"
	"devConnect/config"
	"devConnect/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// We can use the state parameter to pass a redirect URL if needed
	state := r.URL.Query().Get("redirect_url")
	url := config.GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(res http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	if code == "" {
		http.Error(res, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Token Exchange Error:", err)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	oauth2Service, err := oauth2.NewService(context.Background(), option.WithTokenSource(config.GoogleOauthConfig.TokenSource(context.Background(), token)))
	if err != nil {
		log.Println("OAuth2 Service Error:", err)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Println("User Info Error:", err)
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
		userInfo.Id,
		userInfo.Name,
		userInfo.Email,
		userInfo.Picture,
		time.Now(),
	)

	if err != nil {
		log.Println("Insert error:", err)
	}

	// Generate JWT Token
	jwtToken, err := utils.GenerateToken(userInfo.Id)
	if err != nil {
		log.Println("JWT Generation Error:", err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var targetURL string
	if state != "" && state != "state" {
		// Use state as redirect URL (make sure to handle URL params)
		separator := "?"
		if strings.Contains(state, "?") {
			separator = "&"
		}
		targetURL = fmt.Sprintf("%s%suserId=%s&name=%s&token=%s&success=true", state, separator, userInfo.Id, userInfo.Name, jwtToken)
	} else {
		// Default to mobile deep link
		targetURL = fmt.Sprintf("devconnect://auth?userId=%s&name=%s&token=%s&success=true", userInfo.Id, userInfo.Name, jwtToken)
	}

	res.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(res, `
		<html>
			<body style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; font-family: sans-serif; text-align: center; padding: 20px;">
				<h2 style="color: #4285F4;">Login Successful!</h2>
				<p>You are being redirected back to the DevConnect app.</p>
				<p style="font-size: 0.9em; color: #666;">If you are not redirected automatically, click the button below:</p>
				<a href="%s" style="display: inline-block; padding: 12px 24px; background: #4285F4; color: white; text-decoration: none; border-radius: 8px; margin-top: 20px; font-weight: bold; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">Return to App</a>
				<script>
					window.location.href = "%s";
					setTimeout(function() {
						window.location.href = "%s";
					}, 2000);
				</script>
			</body>
		</html>
	`, targetURL, targetURL, targetURL)
}
