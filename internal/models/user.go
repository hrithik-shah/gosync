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

	RootDirectoryID uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`

	RootDirectory *Directory `gorm:"foreignKey:RootDirectoryID;constraint:fk_users_root_directory,deferrable:initially deferred"`

	Devices     []Device
	Directories []Directory
	Files       []File
	Events      []Event
}

func (User) ForeignKeys() []string {
	return []string{"RootDirectory"}
}
