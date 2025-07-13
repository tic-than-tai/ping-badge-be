package model

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ActivityID   uuid.UUID  `json:"activity_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgID        uuid.UUID  `json:"org_id" gorm:"type:uuid;not null"`
	ActivityName string     `json:"activity_name" gorm:"type:varchar(255);not null"`
	Description  *string    `json:"description" gorm:"type:text"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Location     *string    `json:"location" gorm:"type:varchar(255)"`
	BadgeDefID   *uuid.UUID `json:"badge_def_id" gorm:"type:uuid"`
	BaseModel

	// Relationships
	Organization   Organization            `gorm:"-"`
	Badge          *Badge                  `gorm:"-"`
	Participations []ActivityParticipation `gorm:"-"`
}
