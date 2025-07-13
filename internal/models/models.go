package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model with common fields
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// User represents the users table
type User struct {
	UserID             uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username           string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email              string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash       string    `json:"-" gorm:"type:varchar(255);not null"`
	FullName           *string   `json:"full_name" gorm:"type:varchar(100)"`
	ProfilePictureURL  *string   `json:"profile_picture_url" gorm:"type:varchar(255)"`
	Bio                *string   `json:"bio" gorm:"type:text"`
	Role               string    `json:"role" gorm:"type:varchar(20);default:'USER'"`
	PrivacySetting     string    `json:"privacy_setting" gorm:"type:varchar(20);default:'public'"`
	BaseModel
}

// Organization represents the organizations table
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
	Owner  User                   `gorm:"-"`
	Admins []OrganizationAdmin    `gorm:"-"`
	Badges []Badge                `gorm:"-"`
}

// OrganizationAdmin represents the organization_admins table
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

// Badge represents the badges table
type Badge struct {
	BadgeDefID  uuid.UUID              `json:"badge_def_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgID       uuid.UUID              `json:"org_id" gorm:"type:uuid;not null;index:idx_org_badge_name,unique"`
	BadgeName   string                 `json:"badge_name" gorm:"type:varchar(100);not null;index:idx_org_badge_name,unique"`
	Description *string                `json:"description" gorm:"type:text"`
	ImageURL    string                 `json:"image_url" gorm:"type:varchar(255);not null"`
	Criteria    *string                `json:"criteria" gorm:"type:text"`
	BadgeType   string                 `json:"badge_type" gorm:"type:varchar(20);not null;default:'instant'"`
	RuleConfig  map[string]interface{} `json:"rule_config" gorm:"type:jsonb"`
	IsActive    bool                   `json:"is_active" gorm:"default:true"`
	BaseModel

	// Relationships
	Organization Organization     `gorm:"-"`
	IssuedBadges []IssuedBadge    `gorm:"-"`
	Activities   []Activity       `gorm:"-"`
	Progress     []UserBadgeProgress `gorm:"-"`
}

// IssuedBadge represents the issued_badges table
type IssuedBadge struct {
	IssuedBadgeID                uuid.UUID              `json:"issued_badge_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BadgeDefID                   uuid.UUID              `json:"badge_def_id" gorm:"type:uuid;not null;index:idx_user_badge;index:idx_org_badge"`
	UserID                       uuid.UUID              `json:"user_id" gorm:"type:uuid;not null;index:idx_user_badge"`
	OrgID                        uuid.UUID              `json:"org_id" gorm:"type:uuid;not null;index:idx_org_badge"`
	IssueDate                    time.Time              `json:"issue_date" gorm:"autoCreateTime"`
	VerificationCode             string                 `json:"verification_code" gorm:"type:varchar(255);uniqueIndex;not null"`
	SourceType                   *string                `json:"source_type" gorm:"type:varchar(50)"`
	SourceID                     *uuid.UUID             `json:"source_id" gorm:"type:uuid"`
	CumulativeProgressAtIssuance *float64               `json:"cumulative_progress_at_issuance" gorm:"type:numeric"`
	CumulativeUnit               *string                `json:"cumulative_unit" gorm:"type:varchar(50)"`
	AdditionalData               map[string]interface{} `json:"additional_data" gorm:"type:jsonb"`
	Status                       string                 `json:"status" gorm:"type:varchar(20);default:'issued'"`
	BlockchainTxID               *string                `json:"blockchain_tx_id" gorm:"type:varchar(255)"`

	// Relationships
	Badge        Badge        `gorm:"-"`
	User         User         `gorm:"-"`
	Organization Organization `gorm:"-"`
	Views        []BadgeView  `gorm:"-"`
}

// UserBadgeProgress represents the user_badge_progress table
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

// Activity represents the activities table
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

// ActivityParticipation represents the activity_participation table
type ActivityParticipation struct {
	ParticipationID         uuid.UUID    `json:"participation_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ActivityID              uuid.UUID    `json:"activity_id" gorm:"type:uuid;not null"`
	UserID                  uuid.UUID    `json:"user_id" gorm:"type:uuid;not null"`
	Status                  string       `json:"status" gorm:"type:varchar(50);default:'registered'"`
	ProofOfParticipationURL *string      `json:"proof_of_participation_url" gorm:"type:varchar(255)"`
	IssuedBadgeID           *uuid.UUID   `json:"issued_badge_id" gorm:"type:uuid;uniqueIndex"`
	CreatedAt               time.Time    `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Activity    Activity     `gorm:"-"`
	User        User         `gorm:"-"`
	IssuedBadge *IssuedBadge `gorm:"-"`
}

// BadgeView represents the badge_views table
type BadgeView struct {
	ViewID            uuid.UUID `json:"view_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IssuedBadgeID     uuid.UUID `json:"issued_badge_id" gorm:"type:uuid;not null"`
	ViewerIPAddress   *string   `json:"viewer_ip_address" gorm:"type:varchar(45)"`
	ViewTimestamp     time.Time `json:"view_timestamp" gorm:"autoCreateTime"`

	// Relationships
	IssuedBadge IssuedBadge `gorm:"-"`
}
