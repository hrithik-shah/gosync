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

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`

	DeletedAt *time.Time

	User User `gorm:"foreignKey:UserID;constraint:fk_files_user,deferrable:initially deferred"`

	Directory Directory `gorm:"foreignKey:DirectoryID;constraint:fk_files_directory,deferrable:initially deferred"`

	CurrentVersion *FileVersion `gorm:"foreignKey:CurrentVersionID;constraint:fk_files_current_version,deferrable:initially deferred"`

	Versions []FileVersion
}

func (File) ForeignKeys() []string {
	return []string{"User", "Directory", "CurrentVersion"}
}
