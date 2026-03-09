package routes

import (
	"devConnect/handlers"
	"devConnect/middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func SetupRoutes() *mux.Router {

	secret := os.Getenv("SESSION_SECRET")

	// Create session store
	store := sessions.NewCookieStore([]byte(secret))

	// Important for Safari + OAuth
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	}

	gothic.Store = store

	// Tell goth how to extract provider from URL
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return mux.Vars(req)["provider"], nil
	}

	// Configure Google OAuth
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			"http://localhost:3000/auth/google/callback",
		),
	)

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/", homeHandler).Methods("GET")

	r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler).Methods("GET")

	r.HandleFunc("/auth/{provider}/callback", handlers.GoogleCallback).Methods("GET")

	// Protected routes
	r.Handle("/me",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.GetCurrentUser),
		)).Methods("GET")
	r.Handle("/posts",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.CreatePost),
		)).Methods("POST")

	r.Handle("/posts",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.GetPosts),
		)).Methods("GET")

	r.Handle("/follow/{userID}",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.FollowUser),
		)).Methods("POST")

	r.HandleFunc("/followers/{userID}", handlers.GetFollowers).Methods("GET")

	r.HandleFunc("/following/{userID}", handlers.GetFollowing).Methods("GET")

	r.Handle("/feed",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.GetFeed),
		)).Methods("GET")
	r.Handle("/profile",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.UpdateProfile),
		)).Methods("PUT")

	r.HandleFunc("/users/{userID}", handlers.GetUserProfile).Methods("GET")

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
