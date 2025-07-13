package handlers

import (
	"net/http"
	"strconv"

	"ping-badge-be/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityHandler struct {
	db *gorm.DB
}

func NewActivityHandler(db *gorm.DB) *ActivityHandler {
	return &ActivityHandler{db: db}
}

type CreateActivityRequest struct {
	ActivityName string `json:"activity_name" binding:"required"`
	Description  string `json:"description"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Location     string `json:"location"`
	BadgeDefID   string `json:"badge_def_id"`
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner or admin can create activities"})
		return
	}

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity := models.Activity{
		ActivityID:   uuid.New(),
		OrgID:        uuid.MustParse(orgID),
		ActivityName: req.ActivityName,
	}

	if req.Description != "" {
		activity.Description = &req.Description
	}
	if req.Location != "" {
		activity.Location = &req.Location
	}
	if req.BadgeDefID != "" {
		if badgeUUID, err := uuid.Parse(req.BadgeDefID); err == nil {
			activity.BadgeDefID = &badgeUUID
		}
	}

	// Parse dates if provided
	if req.StartDate != "" {
		// You might want to add proper date parsing here
		// For now, we'll skip date parsing and let the database handle it
	}
	if req.EndDate != "" {
		// You might want to add proper date parsing here
	}

	if err := h.db.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	// Load activity with relationships
	if err := h.db.Preload("Organization").Preload("Badge").First(&activity, activity.ActivityID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load activity details"})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	orgID := c.Query("org_id")
	var activities []models.Activity
	query := h.db.Preload("Organization").Preload("Badge")

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

	if err := query.Offset(offset).Limit(limitInt).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (h *ActivityHandler) GetActivity(c *gin.Context) {
	activityID := c.Param("id")
	if _, err := uuid.Parse(activityID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	var activity models.Activity
	if err := h.db.Preload("Organization").Preload("Badge").Preload("Participations.User").First(&activity, "activity_id = ?", activityID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	activityID := c.Param("id")
	if _, err := uuid.Parse(activityID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var activity models.Activity
	if err := h.db.Preload("Organization").First(&activity, "activity_id = ?", activityID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity"})
		return
	}

	// Check if user is owner or admin of the organization
	isAuthorized := activity.Organization.UserIDOwner == userID.(uuid.UUID)
	if !isAuthorized {
		var admin models.OrganizationAdmin
		if err := h.db.Where("org_id = ? AND user_id = ?", activity.OrgID, userID).First(&admin).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner or admin can update activities"})
		return
	}

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	updates := map[string]interface{}{
		"activity_name": req.ActivityName,
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.BadgeDefID != "" {
		if badgeUUID, err := uuid.Parse(req.BadgeDefID); err == nil {
			updates["badge_def_id"] = badgeUUID
		}
	}

	if err := h.db.Model(&activity).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}

	// Load updated activity
	if err := h.db.Preload("Organization").Preload("Badge").First(&activity, "activity_id = ?", activityID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated activity"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	activityID := c.Param("id")
	if _, err := uuid.Parse(activityID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var activity models.Activity
	if err := h.db.Preload("Organization").First(&activity, "activity_id = ?", activityID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity"})
		return
	}

	// Check if user is owner or admin of the organization
	isAuthorized := activity.Organization.UserIDOwner == userID.(uuid.UUID)
	if !isAuthorized {
		var admin models.OrganizationAdmin
		if err := h.db.Where("org_id = ? AND user_id = ?", activity.OrgID, userID).First(&admin).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner or admin can delete activities"})
		return
	}

	if err := h.db.Delete(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

// Join an activity
func (h *ActivityHandler) JoinActivity(c *gin.Context) {
	activityID := c.Param("id")
	if _, err := uuid.Parse(activityID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if activity exists
	var activity models.Activity
	if err := h.db.First(&activity, "activity_id = ?", activityID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity"})
		return
	}

	// Check if user is already participating
	var existingParticipation models.ActivityParticipation
	if err := h.db.Where("activity_id = ? AND user_id = ?", activityID, userID).First(&existingParticipation).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already participating in this activity"})
		return
	}

	// Create participation record
	participation := models.ActivityParticipation{
		ParticipationID: uuid.New(),
		ActivityID:      uuid.MustParse(activityID),
		UserID:          userID.(uuid.UUID),
		Status:          "registered",
	}

	if err := h.db.Create(&participation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join activity"})
		return
	}

	// Load participation with relationships
	if err := h.db.Preload("Activity").Preload("User").First(&participation, participation.ParticipationID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load participation details"})
		return
	}

	c.JSON(http.StatusCreated, participation)
}

// Get user's activity participations
func (h *ActivityHandler) GetUserParticipations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var participations []models.ActivityParticipation
	if err := h.db.Preload("Activity.Organization").Preload("Activity.Badge").Where("user_id = ?", userID).Find(&participations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch participations"})
		return
	}

	c.JSON(http.StatusOK, participations)
}
