package main

import (
	"log"

	"gosync/internal/config"
	"gosync/internal/database"
	"gosync/internal/models"

	"gorm.io/gorm"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")

}

func migrate() error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Device{},
		&models.Directory{},
		&models.File{},
		&models.FileVersion{},
		&models.Event{},
		&models.RefreshToken{},
	); err != nil {
		return err
	}

	addConstraints(db)

	return nil
}

func addConstraints(db *gorm.DB) error {
	// Add foreign key constraints for User and Directory
	// Use DEFERRABLE INITIALLY DEFERRED to allow creating
	// User and its root Directory in the same transaction.
	if err := db.Exec(`
		ALTER TABLE directories
		ADD CONSTRAINT fk_directory_user
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		DEFERRABLE INITIALLY DEFERRED
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		ALTER TABLE users
		ADD CONSTRAINT fk_user_root_directory
		FOREIGN KEY (root_directory_id)
		REFERENCES directories(id)
		DEFERRABLE INITIALLY DEFERRED
	`).Error; err != nil {
		return err
	}

	// Enforce no two active files share a name within the same directory
	// — the real, race-safe guard against concurrent duplicate uploads;
	// application-level pre-checks alone can't prevent this under
	// concurrency.
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_file_name_per_directory
		ON files (directory_id, name)
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	return nil
}
