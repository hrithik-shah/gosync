package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"not null;default:'USER'"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Devices     []Device
	Directories []Directory
	Files       []File
	Events      []SyncEvent
}