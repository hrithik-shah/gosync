package main

import (
	"log"
	"net/http"
	"os"

	"gosync/internal/api/repository"
	"gosync/internal/api/router"
	"gosync/internal/config"
	"gosync/internal/database"
)

// @title           GoSync API
// @version         1.0
// @description     File sync service API.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cmd := "start"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	// Load configuration and environment variables
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Set the default database connection for the repository package
	repository.SetDefault(db)

	switch cmd {
	case "start":
		startServer()
	case "migrate":
		if err := database.Migrate(); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		log.Println("migrations applied")
	default:
		log.Fatalf("unknown command: %q (expected \"start\" or \"migrate\")", cmd)
	}
}

func startServer() {
	r := router.New()

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
