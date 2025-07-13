package api_impl

import (
	"context"
	"net/http"
	"strconv"

	"ping-badge-be/internal/model"
	"ping-badge-be/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BadgeAPI struct {
	service *service.BadgeService
}

func NewBadgeAPI(service *service.BadgeService) *BadgeAPI {
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

func (api *BadgeAPI) CreateBadge(c *gin.Context) {
	orgID := c.Param("org_id")
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req CreateBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	badge := &model.Badge{
		BadgeDefID: uuid.New(),
		OrgID:      uuid.MustParse(orgID),
		BadgeName:  req.BadgeName,
		ImageURL:   req.ImageURL,
		BadgeType:  req.BadgeType,
		RuleConfig: req.RuleConfig,
		IsActive:   true,
	}
	if req.Description != "" {
		badge.Description = &req.Description
	}
	if req.Criteria != "" {
		badge.Criteria = &req.Criteria
	}

	if err := api.service.CreateBadge(context.Background(), badge); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create badge"})
		return
	}
	c.JSON(http.StatusCreated, badge)
}

func (api *BadgeAPI) GetBadges(c *gin.Context) {
	orgID := c.Query("org_id")
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
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	if pageInt < 1 {
		pageInt = 1
	}
	if limitInt < 1 || limitInt > 100 {
		limitInt = 10
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
	badge.ImageURL = req.ImageURL
	badge.BadgeType = req.BadgeType
	badge.RuleConfig = req.RuleConfig
	if req.Description != "" {
		badge.Description = &req.Description
	}
	if req.Criteria != "" {
		badge.Criteria = &req.Criteria
	}
	if err := api.service.UpdateBadge(context.Background(), badge); err != nil {
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
	if err := api.service.DeleteBadge(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete badge"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Badge deleted successfully"})
}
