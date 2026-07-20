package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	DirectoryID uuid.UUID `gorm:"type:uuid;not null;index"`

	Name string `gorm:"not null"`

	// Latest version
	CurrentVersionID *uuid.UUID `gorm:"type:uuid"`

	// Hash of current content
	ContentHash string

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt *time.Time

	User User `gorm:"foreignKey:UserID"`

	Directory Directory `gorm:"foreignKey:DirectoryID"`

	CurrentVersion *FileVersion `gorm:"foreignKey:CurrentVersionID"`

	Versions []FileVersion
}
