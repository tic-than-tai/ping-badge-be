package api_impl

import (
	"context"
	"net/http"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BadgeAPI struct {
	service service.BadgeService
}

func NewBadgeAPI(service service.BadgeService) *BadgeAPI {
	return &BadgeAPI{service: service}
}

type CreateBadgeRequest struct {
	BadgeName   string                 `json:"badge_name" binding:"required"`
	Description string                 `json:"description"`
	ImageURL    string                 `json:"image_url" binding:"required"`
	Criteria    string                 `json:"criteria"`
	BadgeType   string                 `json:"badge_type" binding:"required"`
	RuleConfig  map[string]interface{} `json:"rule_config"`
}

func (api *BadgeAPI) ListBadges(c *gin.Context) {
	orgID := c.Query("org_id")
	userID := c.Query("user_id")

	// If user_id is provided, return issued badges for that user
	if userID != "" {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		issuedBadges, err := api.service.ListIssuedBadgesByUser(context.Background(), userUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user badges"})
			return
		}
		c.JSON(http.StatusOK, issuedBadges)
		return
	}

	// Original logic for listing all badges by organization
	var orgUUID *uuid.UUID
	if orgID != "" {
		parsed, err := uuid.Parse(orgID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}
		orgUUID = &parsed
	}
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	pageInt := 1
	limitInt := 10
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageInt = p
	}
	if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
		limitInt = l
	}
	offset := (pageInt - 1) * limitInt
	badges, err := api.service.ListBadges(context.Background(), orgUUID, offset, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badges"})
		return
	}
	c.JSON(http.StatusOK, badges)
}

func (api *BadgeAPI) GetBadge(c *gin.Context) {
	badgeID := c.Param("id")
	id, err := uuid.Parse(badgeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
		return
	}
	badge, err := api.service.GetBadge(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
		return
	}
	c.JSON(http.StatusOK, badge)
}

func (api *BadgeAPI) CreateBadge(c *gin.Context) {
	orgID := c.Param("id")
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	var req CreateBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	badge := &model.Badge{
		BadgeDefID:  uuid.New(),
		OrgID:       orgUUID,
		BadgeName:   req.BadgeName,
		Description: &req.Description,
		ImageURL:    req.ImageURL,
		Criteria:    &req.Criteria,
		BadgeType:   req.BadgeType,
		RuleConfig:  req.RuleConfig,
		IsActive:    true,
	}
	err = api.service.CreateBadge(context.Background(), badge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create badge"})
		return
	}
	c.JSON(http.StatusCreated, badge)
}

func (api *BadgeAPI) UpdateBadge(c *gin.Context) {
	badgeID := c.Param("id")
	id, err := uuid.Parse(badgeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
		return
	}
	var req CreateBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	badge, err := api.service.GetBadge(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
		return
	}
	badge.BadgeName = req.BadgeName
	badge.Description = &req.Description
	badge.ImageURL = req.ImageURL
	badge.Criteria = &req.Criteria
	badge.BadgeType = req.BadgeType
	badge.RuleConfig = req.RuleConfig
	err = api.service.UpdateBadge(context.Background(), badge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update badge"})
		return
	}
	c.JSON(http.StatusOK, badge)
}

func (api *BadgeAPI) DeleteBadge(c *gin.Context) {
	badgeID := c.Param("id")
	id, err := uuid.Parse(badgeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
		return
	}
	err = api.service.DeleteBadge(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete badge"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Badge deleted successfully"})
}
