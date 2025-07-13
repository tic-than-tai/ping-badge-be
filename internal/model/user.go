package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID            uuid.UUID
	Username          string
	Email             string
	PasswordHash      string
	FullName          *string
	ProfilePictureURL *string
	Bio               *string
	Role              string
	PrivacySetting    string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
