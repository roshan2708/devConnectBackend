package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"devConnect/config"
	"devConnect/middleware"
	"devConnect/routes"
	"devConnect/workers"
)

func main() {

	// Load env
	config.LoadEnv()

	// Initialize OAuth
	config.InitOauth()

	// Connect DB
	config.ConnectDB()

	// Start Background Workers
	workers.StartSyncWorker()

	// Setup router
	r := routes.SetupRoutes()

	// Get Render port
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	fmt.Println("Server running on port", port)
	fmt.Println("http://localhost:3000")

	log.Fatal(http.ListenAndServe(":"+port, middleware.CorsMiddleware(r)))
}
