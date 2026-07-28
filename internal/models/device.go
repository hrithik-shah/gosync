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

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`

	User User `gorm:"foreignKey:UserID;constraint:fk_devices_user,deferrable:initially deferred"`
}

func (Device) ForeignKeys() []string {
	return []string{"User"}
}
