package models

import (
	"time"

	"gorm.io/gorm"
)

// UserTokens represents the total tokens for a user in PostgreSQL.
type UserTokens struct {
	gorm.Model
	UserID      string `gorm:"uniqueIndex"`
	Domain      string `gorm:"uniqueIndex"`
	TotalTokens int32
}

// UsedToken represents the token usage records for a user in PostgreSQL.
type UsedToken struct {
	gorm.Model
	UserID string `gorm:"uniqueIndex"`
	Domain string `gorm:"uniqueIndex"`
	Tokens int32
	UsedAt time.Time
}
