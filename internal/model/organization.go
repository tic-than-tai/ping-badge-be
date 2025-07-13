package model

import (
	"github.com/google/uuid"
)

type Organization struct {
	OrgID       uuid.UUID `json:"org_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgName     string    `json:"org_name" gorm:"type:varchar(255);uniqueIndex;not null"`
	OrgEmail    string    `json:"org_email" gorm:"type:varchar(100);uniqueIndex;not null"`
	OrgLogoURL  *string   `json:"org_logo_url" gorm:"type:varchar(255)"`
	UserIDOwner uuid.UUID `json:"user_id_owner" gorm:"type:uuid"`
	Description *string   `json:"description" gorm:"type:text"`
	WebsiteURL  *string   `json:"website_url" gorm:"type:varchar(255)"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	BaseModel

	// Relationships
	Owner  User                `gorm:"-"`
	Admins []OrganizationAdmin `gorm:"-"`
	Badges []Badge             `gorm:"-"`
}
