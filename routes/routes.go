package routes

import (
	"devConnect/handlers"
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

	store := sessions.NewCookieStore([]byte(secret))

	gothic.Store = store

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			"http://localhost:3000/auth/google/callback",
		),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)

	r.HandleFunc("/auth/{provider}/callback", handlers.GoogleCallback)
	r.HandleFunc("/me", handlers.GetCurrentUser).Methods("GET")

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
