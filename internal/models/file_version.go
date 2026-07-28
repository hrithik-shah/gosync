package models

import (
	"time"

	"github.com/google/uuid"
)

type FileVersion struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	FileID uuid.UUID `gorm:"type:uuid;not null"`

	VersionNumber int `gorm:"not null"`

	// Path/key in filesystem or MinIO
	StorageKey string `gorm:"not null"`

	SizeBytes int64 `gorm:"not null"`

	ContentHash string `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`

	File File `gorm:"foreignKey:FileID;constraint:fk_file_versions_file,deferrable:initially deferred"`
}

func (FileVersion) ForeignKeys() []string {
	return []string{"File"}
}
