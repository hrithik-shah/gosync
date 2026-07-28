package models

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
	EventTypeMove   EventType = "move"
)

func (t EventType) IsValid() bool {
	switch t {
	case EventTypeCreate, EventTypeUpdate, EventTypeDelete, EventTypeMove:
		return true
	default:
		return false
	}
}

type Event struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	Type EventType `gorm:"type:event_type;not null"`

	FileID      *uuid.UUID `gorm:"type:uuid"`
	DirectoryID *uuid.UUID `gorm:"type:uuid"`

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`

	User      User       `gorm:"foreignKey:UserID;constraint:fk_events_user,deferrable:initially deferred"`
	File      *File      `gorm:"foreignKey:FileID;constraint:fk_events_file,deferrable:initially deferred"`
	Directory *Directory `gorm:"foreignKey:DirectoryID;constraint:fk_events_directory,deferrable:initially deferred"`
}

func (Event) ForeignKeys() []string {
	return []string{"User", "File", "Directory"}
}
