package main

import (
	"log"
	"os"

	"gosync/internal/config"
	"gosync/internal/database"
)

func main() {
	cmd := "start"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case "start":
		startServer(db)
	case "migrate":
		if err := database.Migrate(db); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		log.Println("migrations applied")
	default:
		log.Fatalf("unknown command: %q (expected \"start\" or \"migrate\")", cmd)
	}
}

func startServer(db interface { /* replace with *gorm.DB */
}) {
	// TODO: wire up your actual server (router, listener, etc.)
	log.Println("starting server...")
}
