package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gosync/internal/config"
	"gosync/internal/models"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.POSTGRES_HOST, cfg.POSTGRES_PORT, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_DB, cfg.POSTGRES_SSLMODE,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Device{},
		&models.Directory{},
		&models.File{},
		&models.FileVersion{},
		&models.SyncEvent{},
	)
}
