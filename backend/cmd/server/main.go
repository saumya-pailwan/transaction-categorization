package main

import (
	"log"
	"net/http"
	"os"

	"autonomoustx/internal/api"
	"autonomoustx/internal/db"
	"autonomoustx/internal/plaid"
)

func main() {
	// 1. Connect to DB
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Pool.Close()

	// 2. Run Migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 3. Init Plaid
	plaid.Init()

	// 4. Setup Router
	r := api.NewRouter()

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
