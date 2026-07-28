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

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`

	DeletedAt *time.Time

	User User `gorm:"foreignKey:UserID;constraint:fk_directories_user,deferrable:initially deferred"`

	Parent   *Directory  `gorm:"foreignKey:ParentID;constraint:fk_directories_parent,deferrable:initially deferred"`
	Children []Directory `gorm:"foreignKey:ParentID"`

	Files []File
}

func (Directory) ForeignKeys() []string {
	return []string{"User", "Parent"}
}
