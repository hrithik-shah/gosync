package main

import (
	"log"
	"net/http"
	"os"

	"gosync/internal/config"
	"gosync/internal/database"

	"gorm.io/gorm"
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

func startServer(db *gorm.DB) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
