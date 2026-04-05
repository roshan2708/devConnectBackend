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

	// Construct the target URL for the app
	var targetURL string
	if state != "" && state != "state" {
		separator := "?"
		if strings.Contains(state, "?") {
			separator = "&"
		}
		targetURL = fmt.Sprintf("%s%suserId=%s&name=%s&token=%s&success=true", state, separator, userInfo.Id, userInfo.Name, jwtToken)
	} else {
		targetURL = fmt.Sprintf("devconnect://auth/callback?userId=%s&name=%s&token=%s&success=true", userInfo.Id, userInfo.Name, jwtToken)
	}

	// For Android, intent:// is more robust for opening apps from browsers
	// Our URL: devconnect://auth/callback?...
	intentURL := fmt.Sprintf("intent://auth/callback?userId=%s&name=%s&token=%s&success=true#Intent;scheme=devconnect;package=com.example.dev_connect_app;S.userId=%s;S.token=%s;end", 
		userInfo.Id, userInfo.Name, jwtToken, userInfo.Id, jwtToken)

	res.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(res, `
		<html>
			<head>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Login Successful</title>
				<style>
					body { display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 100vh; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; text-align: center; padding: 20px; background-color: #f8f9fa; margin: 0; }
					.card { background: white; padding: 40px; border-radius: 16px; box-shadow: 0 10px 25px rgba(0,0,0,0.05); max-width: 400px; width: 100%%; }
					.icon-box { background: #e7f4ff; width: 64px; height: 64px; border-radius: 50%%; display: flex; align-items: center; justify-content: center; margin: 0 auto 24px; }
					h2 { color: #1a1f36; margin: 0 0 12px; font-size: 24px; }
					p { color: #4f566b; margin: 0 0 32px; line-height: 1.5; }
					.btn-primary { display: block; padding: 16px 24px; background: #4285F4; color: white !important; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; transition: background 0.2s; box-shadow: 0 4px 6px rgba(66, 133, 244, 0.2); margin-bottom: 12px; border: none; cursor: pointer; width: 100%%; }
					.btn-primary:hover { background: #3367d6; }
					.btn-secondary { color: #4285F4; text-decoration: none; font-size: 0.9em; font-weight: 500; cursor: pointer; }
					.debug-section { margin-top: 40px; padding-top: 20px; border-top: 1px solid #edf2f7; text-align: left; }
					.debug-title { font-size: 12px; color: #a3acb9; margin-bottom: 8px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; }
					code { display: block; font-size: 10px; color: #718096; word-break: break-all; background: #f7fafc; padding: 10px; border-radius: 6px; border: 1px solid #edf2f7; white-space: pre-wrap; }
				</style>
			</head>
			<body>
				<div class="card">
					<div class="icon-box">
						<svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#4285F4" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
					</div>
					<h2>Login Successful!</h2>
					<p>You are being redirected back to the DevConnect app.</p>
					
					<button onclick="tryRedirect()" class="btn-primary">Return to App</button>
					<a href="%s" class="btn-secondary">Use Custom Scheme Redirect</a>
					
					<div class="debug-section">
						<p class="debug-title">Debug Info</p>
						<code>
State: %s
Target: %s
Intent: %s
User Agent: <script>document.write(navigator.userAgent)</script>
						</code>
					</div>
				</div>
				<script>
					const intentURL = "%s";
					const schemeURL = "%s";
					
					function tryRedirect() {
						const isAndroid = /Android/i.test(navigator.userAgent);
						// For Android Chrome, intent:// is more robust.
						// However, let's try opening the app via custom scheme first as it works more universally if configured.
						
						if (isAndroid) {
							// Try intent first
							window.location.href = intentURL;
							
							// Fallback after short delay
							setTimeout(() => {
								window.location.href = schemeURL;
							}, 1500);
						} else {
							window.location.href = schemeURL;
						}
					}

					// Auto-try
					window.onload = function() {
						setTimeout(tryRedirect, 500);
					}
				</script>
			</body>
		</html>
	`, targetURL, state, targetURL, intentURL, intentURL, targetURL)
}
