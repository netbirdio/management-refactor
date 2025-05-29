package team

import (
	"management/internal/shared/db"
	"time"
)

// PersonalAccessToken holds all information about a PAT including a hashed version of it for verification
type PersonalAccessToken struct {
	ID string `gorm:"primaryKey"`
	// User is a reference to Account that this object belongs
	UserID         string `gorm:"index"`
	Name           string
	HashedToken    string
	ExpirationDate *time.Time
	// scope could be added in future
	CreatedBy string
	CreatedAt time.Time
	LastUsed  *time.Time
}

func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}

type PersonalAccessTokenEvent = db.ModelEvent[PersonalAccessToken]
