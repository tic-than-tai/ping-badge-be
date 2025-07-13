package model

import (
	"time"

	"github.com/google/uuid"
)

type UserBadgeProgress struct {
	ProgressID    uuid.UUID `json:"progress_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index:idx_user_badge_progress,unique"`
	BadgeDefID    uuid.UUID `json:"badge_def_id" gorm:"type:uuid;not null;index:idx_user_badge_progress,unique"`
	ProgressValue float64   `json:"progress_value" gorm:"type:numeric;default:0"`
	Unit          *string   `json:"unit" gorm:"type:varchar(50)"`
	IsQualified   bool      `json:"is_qualified" gorm:"default:false"`
	LastUpdated   time.Time `json:"last_updated" gorm:"autoUpdateTime"`

	// Relationships
	User  User  `gorm:"-"`
	Badge Badge `gorm:"-"`
}
