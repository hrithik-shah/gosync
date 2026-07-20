package models

import (
	"time"

	"github.com/google/uuid"
)

type Directory struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	// Parent directory
	// NULL means root directory
	ParentID *uuid.UUID `gorm:"type:uuid"`

	Name string `gorm:"not null"`

	// Hash of children
	MerkleHash string

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt *time.Time

	User User `gorm:"foreignKey:UserID"`

	Parent   *Directory  `gorm:"foreignKey:ParentID"`
	Children []Directory `gorm:"foreignKey:ParentID"`

	Files []File
}
