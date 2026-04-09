package routes

import (
	"devConnect/handlers"
	"devConnect/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.RateLimiter)

	// Public routes
	r.HandleFunc("/", homeHandler).Methods("GET")

	// Manual OAuth2 Routes
	r.HandleFunc("/auth/google", handlers.GoogleLogin).Methods("GET")
	r.HandleFunc("/auth/google/callback", handlers.GoogleCallback).Methods("GET")

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


	// --- Advanced Architecture Extensions ---

	// 1. Platform Links & Gamification
	r.Handle("/integrations/{platform}/link",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.LinkPlatform),
		)).Methods("POST")

	r.HandleFunc("/users/{userID}/stats", handlers.GetUserStats).Methods("GET")

	// 2. DNA & AI
	r.Handle("/users/me/analyze-dna",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.AnalyzeDNA),
		)).Methods("POST")

	r.HandleFunc("/users/{userID}/dna", handlers.GetDNA).Methods("GET")

	// 3. New Advanced Feed & Extended Posts
	r.Handle("/feed/live",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.GetAdvancedFeed),
		)).Methods("GET")

	r.HandleFunc("/posts/advanced", handlers.GetExtendedPosts).Methods("GET")

	// 4. Live Sessions (Build With Me)
	r.Handle("/sessions",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.StartLiveSession),
		)).Methods("POST")

	r.HandleFunc("/sessions/live", handlers.GetLiveSessions).Methods("GET")

	return r

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
