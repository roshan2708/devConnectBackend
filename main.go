package main

import (
	"fmt"
	"log"
	"net/http"

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

	fmt.Println("Server running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
