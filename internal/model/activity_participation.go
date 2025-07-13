package model

import (
	"time"

	"github.com/google/uuid"
)

type ActivityParticipation struct {
	ParticipationID         uuid.UUID  `json:"participation_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ActivityID              uuid.UUID  `json:"activity_id" gorm:"type:uuid;not null"`
	UserID                  uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Status                  string     `json:"status" gorm:"type:varchar(50);default:'registered'"`
	ProofOfParticipationURL *string    `json:"proof_of_participation_url" gorm:"type:varchar(255)"`
	IssuedBadgeID           *uuid.UUID `json:"issued_badge_id" gorm:"type:uuid;uniqueIndex"`
	CreatedAt               time.Time  `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Activity    Activity     `gorm:"-"`
	User        User         `gorm:"-"`
	IssuedBadge *IssuedBadge `gorm:"-"`
}
