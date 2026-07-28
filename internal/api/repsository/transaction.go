package repsository

import (
	"gosync/internal/database"

	"gorm.io/gorm"
)

// WithTransaction runs fn inside a DB transaction. If fn returns an error
// (or panics), the transaction is rolled back; otherwise it's committed.
func WithTransaction(fn func(tx *gorm.DB) error) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
