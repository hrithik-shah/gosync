package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gosync/internal/config"
)

var globalDB *gorm.DB

func Connect() (*gorm.DB, error) {
	cfg := config.Get()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost(), cfg.PostgresPort(), cfg.PostgresUser(), cfg.PostgresPassword(), cfg.PostgresDb(), cfg.PostgresSslmode(),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	globalDB = db

	return db, nil
}

func GetDB() (*gorm.DB, error) {
	if globalDB == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return globalDB, nil
}
