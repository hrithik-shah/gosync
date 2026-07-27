package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"not null;default:'USER'"`

	rootDirectoryID uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	rootDirectory *Directory `gorm:"foreignKey:rootDirectoryID"`

	Devices     []Device
	Directories []Directory
	Files       []File
	Events      []SyncEvent
}
