package model

import (
	"time"

	"github.com/google/uuid"
)

type Badge struct {
	BadgeDefID  uuid.UUID
	OrgID       uuid.UUID
	BadgeName   string
	Description *string
	ImageURL    string
	Criteria    *string
	BadgeType   string
	RuleConfig  map[string]interface{}
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
