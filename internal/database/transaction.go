package database

import "gorm.io/gorm"

// WithTransaction runs fn inside a DB transaction. If fn returns an error
// (or panics), the transaction is rolled back; otherwise it's committed.
func WithTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
