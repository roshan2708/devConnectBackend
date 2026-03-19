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
		SameSite: http.SameSiteLaxMode,
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
	r.Use(middleware.RateLimiter)

	// Public routes
	r.HandleFunc("/", homeHandler).Methods("GET")

	r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler).Methods("GET")

	r.HandleFunc("/auth/{provider}/callback", handlers.GoogleCallback).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

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

	r.Handle("/posts/{postID}/like",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.LikePost),
		)).Methods("POST")

	r.HandleFunc("/posts/{postID}/likes", handlers.GetLikes).Methods("GET")

	r.Handle("/posts/{postID}/comment",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.CreateComment),
		)).Methods("POST")

	r.HandleFunc("/posts/{postID}/comments", handlers.GetComments).Methods("GET")
	r.Handle("/notifications",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.GetNotifications),
		)).Methods("GET")

	r.Handle("/posts/{postID}",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.DeletePost),
		)).Methods("DELETE")

	r.Handle("/posts/{postID}",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.EditPost),
		)).Methods("PUT")

	r.Handle("/posts/{postID}/like",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.UnlikePost),
		)).Methods("DELETE")

	r.Handle("/follow/{userID}",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.UnfollowUser),
		)).Methods("DELETE")

	r.Handle("/comments/{commentID}",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.DeleteComment),
		)).Methods("DELETE")

	r.Handle("/notifications/{id}/read",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.MarkNotificationRead),
		)).Methods("PUT")

	r.HandleFunc("/users/{userID}/posts", handlers.GetUserPosts).Methods("GET")
	r.HandleFunc("/search", handlers.SearchUsers).Methods("GET")
	r.HandleFunc("/trending", handlers.GetTrendingPosts).Methods("GET")

	return r

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
