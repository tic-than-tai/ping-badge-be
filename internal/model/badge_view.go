package model

import (
	"time"

	"github.com/google/uuid"
)

type BadgeView struct {
	ViewID          uuid.UUID `json:"view_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IssuedBadgeID   uuid.UUID `json:"issued_badge_id" gorm:"type:uuid;not null"`
	ViewerIPAddress *string   `json:"viewer_ip_address" gorm:"type:varchar(45)"`
	ViewTimestamp   time.Time `json:"view_timestamp" gorm:"autoCreateTime"`
}
