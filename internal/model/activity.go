package model

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ActivityID   uuid.UUID
	OrgID        uuid.UUID
	ActivityName string
	Description  *string
	StartDate    *time.Time
	EndDate      *time.Time
	Location     *string
	BadgeDefID   *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
