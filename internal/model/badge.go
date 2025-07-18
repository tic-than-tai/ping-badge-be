package model

import (
	"github.com/google/uuid"
)

type Badge struct {
	BadgeDefID  uuid.UUID `json:"badge_def_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgID       uuid.UUID `json:"org_id" gorm:"type:uuid;not null;index:idx_org_badge_name,unique"`
	BadgeName   string    `json:"badge_name" gorm:"type:varchar(100);not null;index:idx_org_badge_name,unique"`
	Description *string   `json:"description" gorm:"type:text"`
	ImageURL    string    `json:"image_url" gorm:"type:varchar(255);not null"`
	Criteria    *string   `json:"criteria" gorm:"type:text"`
	BadgeType   string    `json:"badge_type" gorm:"type:varchar(20);not null;default:'instant'"`
	//RuleConfig  map[string]interface{} `json:"rule_config" gorm:"type:jsonb"`
	IsActive bool `json:"is_active" gorm:"default:true"`
	BaseModel

	// Relationships
	Organization Organization  `gorm:"-"`
	IssuedBadges []IssuedBadge `gorm:"-"`
	Activities   []Activity    `gorm:"-"`
	Progress     []IssuedBadge `gorm:"-"`
}
