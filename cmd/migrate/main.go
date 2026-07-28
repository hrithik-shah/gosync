package main

import (
	"fmt"
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

	if err := createEnumTypes(db); err != nil {
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

	if err := createForeignKeys(db); err != nil {
		return err
	}

	if err := addConstraints(db); err != nil {
		return err
	}

	return nil
}

// createEnumTypes creates Postgres native enum types used by models,
// before AutoMigrate runs — columns referencing these types (e.g.
// Event.Type) depend on the type already existing at table-creation
// time. DO $$ ... $$ makes this idempotent, since Postgres has no
// "CREATE TYPE IF NOT EXISTS".
func createEnumTypes(db *gorm.DB) error {
	return db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
				CREATE TYPE event_type AS ENUM ('create', 'update', 'delete', 'move');
			END IF;
		END$$;
	`).Error
}

// createForeignKeys adds every FK constraint AutoMigrate would normally
// generate, but as an explicit step run afterward — once every table
// already exists, there's no dependency-ordering problem left to solve,
// even for the circular users<->directories pair.
//
// Each model declares its own FK association field names via the
// models.ForeignKeyed interface (see e.g. models/user.go), so this
// function stays generic — adding a new model or relationship only
// requires updating that model's own ForeignKeys() method, not this file.
func createForeignKeys(db *gorm.DB) error {
	m := db.Migrator()

	allModels := []interface{}{
		&models.User{},
		&models.Device{},
		&models.Directory{},
		&models.File{},
		&models.FileVersion{},
		&models.Event{},
		&models.RefreshToken{},
		&models.Event{},
	}

	for _, model := range allModels {
		fk, ok := model.(models.ForeignKeyed)
		if !ok {
			continue // no FK associations
		}

		for _, field := range fk.ForeignKeys() {
			if !m.HasConstraint(model, field) {
				if err := m.CreateConstraint(model, field); err != nil {
					return fmt.Errorf("creating constraint for %T.%s: %w", model, field, err)
				}
				log.Printf("created FK constraint for %T.%s", model, field)
			}
		}
	}

	return nil
}

func addConstraints(db *gorm.DB) error {
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

	// Enforce that an Event refers to exactly one target — a File or a
	// Directory, never both, never neither.
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'chk_event_single_target'
			) THEN
				ALTER TABLE events
				ADD CONSTRAINT chk_event_single_target
				CHECK ((file_id IS NOT NULL)::int + (directory_id IS NOT NULL)::int = 1);
			END IF;
		END$$;
	`).Error; err != nil {
		return err
	}

	return nil
}
