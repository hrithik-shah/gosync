package models

import (
	"time"

	"github.com/google/uuid"
)

type SyncEvent struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	DeviceID *uuid.UUID `gorm:"type:uuid"`

	Operation string `gorm:"not null"`

	ResourceType string `gorm:"not null"`

	ResourceID uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`

	Device *Device `gorm:"foreignKey:DeviceID"`
}
