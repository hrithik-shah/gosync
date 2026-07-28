package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	// SHA-256 hash of the actual refresh token. We never store the raw
	// token — only its hash, so a DB leak alone can't be used to log in.
	TokenHash string `gorm:"not null;uniqueIndex"`

	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time // set on logout/rotation; nil means still valid

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`

	User User `gorm:"foreignKey:UserID;constraint:fk_refresh_tokens_user,deferrable:initially deferred"`
}

func (RefreshToken) ForeignKeys() []string {
	return []string{"User"}
}
