package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`

	Name string `gorm:"not null"`

	LastSyncAt   *time.Time
	LastRootHash string

	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}
