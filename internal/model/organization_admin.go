package model

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationAdmin struct {
	AdminID   uuid.UUID
	OrgID     uuid.UUID
	UserID    uuid.UUID
	Role      string
	CreatedAt time.Time
}
