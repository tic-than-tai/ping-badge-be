package model

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationAdmin struct {
	AdminID   uuid.UUID `json:"admin_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgID     uuid.UUID `json:"org_id" gorm:"type:uuid;not null;index:idx_org_user,unique"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index:idx_org_user,unique"`
	Role      string    `json:"role" gorm:"type:varchar(50);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Organization Organization `gorm:"-"`
	User         User         `gorm:"-"`
}
