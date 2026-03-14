package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"devConnect/config"
	"devConnect/routes"
)

func main() {

	// Load env
	config.LoadEnv()

	// Connect DB
	config.ConnectDB()

	// Setup router
	r := routes.SetupRoutes()

	// Get Render port
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	fmt.Println("Server running on port", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
