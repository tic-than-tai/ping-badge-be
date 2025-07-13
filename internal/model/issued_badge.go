package model

import (
	"time"

	"github.com/google/uuid"
)

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
}
