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
	"github.com/rs/cors"
)

func SetupRoutes() http.Handler {

	secret := os.Getenv("SESSION_SECRET")

	// Create session store
	store := sessions.NewCookieStore([]byte(secret))

	// Cookie settings (for localhost testing)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   false, // IMPORTANT for localhost testing
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
			"https://devconnectbackend-wuej.onrender.com/auth/google/callback",
		),
	)

	r := mux.NewRouter()
	r.Use(middleware.RateLimiter)

	// ---------------- PUBLIC ROUTES ----------------

	r.HandleFunc("/", homeHandler).Methods("GET")

	r.HandleFunc("/auth/{provider}", handlers.BeginGoogleAuth).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handlers.GoogleCallback).Methods("GET")

	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// ---------------- PROTECTED ROUTES ----------------

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

	r.Handle("/follow/{userID}",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.UnfollowUser),
		)).Methods("DELETE")

	r.Handle("/feed",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.GetFeed),
		)).Methods("GET")

	r.Handle("/profile",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.UpdateProfile),
		)).Methods("PUT")

	r.Handle("/notifications",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.GetNotifications),
		)).Methods("GET")

	r.Handle("/notifications/{id}/read",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.MarkNotificationRead),
		)).Methods("PUT")

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
			http.HandlerFunc(handlers.LikePost),
		)).Methods("POST")

	r.Handle("/posts/{postID}/like",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.UnlikePost),
		)).Methods("DELETE")

	r.Handle("/posts/{postID}/comment",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.CreateComment),
		)).Methods("POST")

	r.Handle("/comments/{commentID}",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.DeleteComment),
		)).Methods("DELETE")

	// ---------------- PUBLIC DATA ROUTES ----------------

	r.HandleFunc("/followers/{userID}", handlers.GetFollowers).Methods("GET")
	r.HandleFunc("/following/{userID}", handlers.GetFollowing).Methods("GET")
	r.HandleFunc("/posts/{postID}/likes", handlers.GetLikes).Methods("GET")
	r.HandleFunc("/posts/{postID}/comments", handlers.GetComments).Methods("GET")
	r.HandleFunc("/users/{userID}", handlers.GetUserProfile).Methods("GET")
	r.HandleFunc("/users/{userID}/posts", handlers.GetUserPosts).Methods("GET")
	r.HandleFunc("/search", handlers.SearchUsers).Methods("GET")
	r.HandleFunc("/trending", handlers.GetTrendingPosts).Methods("GET")

	// ---------------- CORS CONFIG ----------------

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8000",
			"https://eaf9f920-41bd-46cc-87f0-5d0368dc54d1.lovableproject.com",
			"https://id-preview--eaf9f920-41bd-46cc-87f0-5d0368dc54d1.lovable.app",
		},

		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization",
		},
		AllowCredentials: true,
	})

	return c.Handler(r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DevConnect API running SHER"))
}
