package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gosync/internal/config"
	"gosync/internal/models"
)

var globalDB *gorm.DB

func Connect() (*gorm.DB, error) {
	cfg := config.Get()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost(), cfg.PostgresPort(), cfg.PostgresUser(), cfg.PostgresPassword(), cfg.PostgresDb(), cfg.PostgresSslmode(),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	globalDB = db

	return db, nil
}

func Migrate() error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.User{},
		&models.Device{},
		&models.Directory{},
		&models.File{},
		&models.FileVersion{},
		&models.SyncEvent{},
		&models.RefreshToken{},
	)
}

func GetDB() (*gorm.DB, error) {
	if globalDB == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return globalDB, nil
}
