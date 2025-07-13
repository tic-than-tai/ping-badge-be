package handlers

import (
	"net/http"
	"strconv"

	"ping-badge-be/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BadgeHandler struct {
	db *gorm.DB
}

func NewBadgeHandler(db *gorm.DB) *BadgeHandler {
	return &BadgeHandler{db: db}
}

type CreateBadgeRequest struct {
	BadgeName   string                 `json:"badge_name" binding:"required"`
	Description string                 `json:"description"`
	ImageURL    string                 `json:"image_url" binding:"required"`
	Criteria    string                 `json:"criteria"`
	BadgeType   string                 `json:"badge_type" binding:"required"`
	RuleConfig  map[string]interface{} `json:"rule_config"`
}

func (h *BadgeHandler) CreateBadge(c *gin.Context) {
	orgID := c.Param("org_id")
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is owner or admin of the organization
	var org models.Organization
	if err := h.db.First(&org, "org_id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization"})
		return
	}

	isAuthorized := org.UserIDOwner == userID.(uuid.UUID)
	if !isAuthorized {
		var admin models.OrganizationAdmin
		if err := h.db.Where("org_id = ? AND user_id = ?", orgID, userID).First(&admin).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner or admin can create badges"})
		return
	}

	var req CreateBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if badge name already exists in this organization
	var existingBadge models.Badge
	if err := h.db.Where("org_id = ? AND badge_name = ?", orgID, req.BadgeName).First(&existingBadge).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Badge name already exists in this organization"})
		return
	}

	badge := models.Badge{
		BadgeDefID:  uuid.New(),
		OrgID:       uuid.MustParse(orgID),
		BadgeName:   req.BadgeName,
		ImageURL:    req.ImageURL,
		BadgeType:   req.BadgeType,
		RuleConfig:  req.RuleConfig,
	}

	if req.Description != "" {
		badge.Description = &req.Description
	}
	if req.Criteria != "" {
		badge.Criteria = &req.Criteria
	}

	if err := h.db.Create(&badge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create badge"})
		return
	}

	// Load badge with organization details
	if err := h.db.Preload("Organization").First(&badge, badge.BadgeDefID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load badge details"})
		return
	}

	c.JSON(http.StatusCreated, badge)
}

func (h *BadgeHandler) GetBadges(c *gin.Context) {
	orgID := c.Query("org_id")
	var badges []models.Badge
	query := h.db.Preload("Organization")

	if orgID != "" {
		if _, err := uuid.Parse(orgID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}
		query = query.Where("org_id = ?", orgID)
	}

	// Add pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	
	var pageInt, limitInt int = 1, 10
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageInt = p
	}
	if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
		limitInt = l
	}

	offset := (pageInt - 1) * limitInt

	if err := query.Offset(offset).Limit(limitInt).Find(&badges).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badges"})
		return
	}

	c.JSON(http.StatusOK, badges)
}

func (h *BadgeHandler) GetBadge(c *gin.Context) {
	badgeID := c.Param("id")
	if _, err := uuid.Parse(badgeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
		return
	}

	var badge models.Badge
	if err := h.db.Preload("Organization").First(&badge, "badge_def_id = ?", badgeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badge"})
		return
	}

	c.JSON(http.StatusOK, badge)
}

func (h *BadgeHandler) UpdateBadge(c *gin.Context) {
	badgeID := c.Param("id")
	if _, err := uuid.Parse(badgeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var badge models.Badge
	if err := h.db.Preload("Organization").First(&badge, "badge_def_id = ?", badgeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badge"})
		return
	}

	// Check if user is owner or admin of the organization
	isAuthorized := badge.Organization.UserIDOwner == userID.(uuid.UUID)
	if !isAuthorized {
		var admin models.OrganizationAdmin
		if err := h.db.Where("org_id = ? AND user_id = ?", badge.OrgID, userID).First(&admin).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner or admin can update badges"})
		return
	}

	var req CreateBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	updates := map[string]interface{}{
		"badge_name": req.BadgeName,
		"image_url":  req.ImageURL,
		"badge_type": req.BadgeType,
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Criteria != "" {
		updates["criteria"] = req.Criteria
	}
	if req.RuleConfig != nil {
		updates["rule_config"] = req.RuleConfig
	}

	if err := h.db.Model(&badge).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update badge"})
		return
	}

	// Load updated badge
	if err := h.db.Preload("Organization").First(&badge, "badge_def_id = ?", badgeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated badge"})
		return
	}

	c.JSON(http.StatusOK, badge)
}

func (h *BadgeHandler) DeleteBadge(c *gin.Context) {
	badgeID := c.Param("id")
	if _, err := uuid.Parse(badgeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var badge models.Badge
	if err := h.db.Preload("Organization").First(&badge, "badge_def_id = ?", badgeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badge"})
		return
	}

	// Check if user is owner or admin of the organization
	isAuthorized := badge.Organization.UserIDOwner == userID.(uuid.UUID)
	if !isAuthorized {
		var admin models.OrganizationAdmin
		if err := h.db.Where("org_id = ? AND user_id = ?", badge.OrgID, userID).First(&admin).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner or admin can delete badges"})
		return
	}

	if err := h.db.Delete(&badge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete badge"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Badge deleted successfully"})
}

// Issue a badge to a user
func (h *BadgeHandler) IssueBadge(c *gin.Context) {
	badgeID := c.Param("id")
	if _, err := uuid.Parse(badgeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
		return
	}

	issuerUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		UserID                       string                 `json:"user_id" binding:"required"`
		SourceType                   string                 `json:"source_type"`
		SourceID                     string                 `json:"source_id"`
		CumulativeProgressAtIssuance *float64               `json:"cumulative_progress_at_issuance"`
		CumulativeUnit               string                 `json:"cumulative_unit"`
		AdditionalData               map[string]interface{} `json:"additional_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipientUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if badge exists
	var badge models.Badge
	if err := h.db.Preload("Organization").First(&badge, "badge_def_id = ?", badgeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badge"})
		return
	}

	// Check if user is authorized to issue this badge
	isAuthorized := badge.Organization.UserIDOwner == issuerUserID.(uuid.UUID)
	if !isAuthorized {
		var admin models.OrganizationAdmin
		if err := h.db.Where("org_id = ? AND user_id = ?", badge.OrgID, issuerUserID).First(&admin).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner or admin can issue badges"})
		return
	}

	// Check if recipient user exists
	var recipientUser models.User
	if err := h.db.First(&recipientUser, "user_id = ?", recipientUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipient user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipient user"})
		return
	}

	// Generate verification code
	verificationCode := uuid.New().String()

	// Create issued badge
	issuedBadge := models.IssuedBadge{
		IssuedBadgeID:    uuid.New(),
		BadgeDefID:       badge.BadgeDefID,
		UserID:           recipientUserID,
		OrgID:            badge.OrgID,
		VerificationCode: verificationCode,
		AdditionalData:   req.AdditionalData,
	}

	if req.SourceType != "" {
		issuedBadge.SourceType = &req.SourceType
	}
	if req.SourceID != "" {
		sourceUUID, err := uuid.Parse(req.SourceID)
		if err == nil {
			issuedBadge.SourceID = &sourceUUID
		}
	}
	if req.CumulativeProgressAtIssuance != nil {
		issuedBadge.CumulativeProgressAtIssuance = req.CumulativeProgressAtIssuance
	}
	if req.CumulativeUnit != "" {
		issuedBadge.CumulativeUnit = &req.CumulativeUnit
	}

	if err := h.db.Create(&issuedBadge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue badge"})
		return
	}

	// Load issued badge with relationships
	if err := h.db.Preload("Badge").Preload("User").Preload("Organization").First(&issuedBadge, issuedBadge.IssuedBadgeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load issued badge details"})
		return
	}

	c.JSON(http.StatusCreated, issuedBadge)
}
