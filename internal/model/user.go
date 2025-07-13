package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type User struct {
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username          string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email             string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash      string    `json:"-" gorm:"type:varchar(255);not null"`
	FullName          *string   `json:"full_name" gorm:"type:varchar(100)"`
	ProfilePictureURL *string   `json:"profile_picture_url" gorm:"type:varchar(255)"`
	Bio               *string   `json:"bio" gorm:"type:text"`
	Role              string    `json:"role" gorm:"type:varchar(20);default:'USER'"`
	PrivacySetting    string    `json:"privacy_setting" gorm:"type:varchar(20);default:'public'"`
	BaseModel
}
