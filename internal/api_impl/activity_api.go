package api_impl

import (
	"context"
	"net/http"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"
	"ping-badge-be/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivityAPI struct {
	service              service.ActivityService
	organizationRepo     *repository.OrganizationRepository
	badgeRepo            repository.BadgeRepository
	participationService service.ActivityParticipationService
}

func NewActivityAPI(
	service service.ActivityService,
	organizationRepo *repository.OrganizationRepository,
	badgeRepo repository.BadgeRepository,
	participationService service.ActivityParticipationService,
) *ActivityAPI {
	return &ActivityAPI{
		service:              service,
		organizationRepo:     organizationRepo,
		badgeRepo:            badgeRepo,
		participationService: participationService,
	}
}

type CreateActivityRequest struct {
	ActivityName string `json:"activity_name" binding:"required"`
	Description  string `json:"description"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Location     string `json:"location"`
	BadgeDefID   string `json:"badge_def_id"`
}

func (api *ActivityAPI) ListActivities(c *gin.Context) {
	orgID := c.Query("org_id")
	userID := c.Query("user_id")

	// If user_id is provided, return activities that the user has participated in
	if userID != "" {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
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

		activities, err := api.service.ListActivitiesByUser(context.Background(), userUUID, offset, limitInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user activities"})
			return
		}
		c.JSON(http.StatusOK, activities)
		return
	}

	// Original logic for listing all activities by organization
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
	activities, err := api.service.ListActivities(context.Background(), orgUUID, offset, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (api *ActivityAPI) GetActivity(c *gin.Context) {
	activityID := c.Param("id")
	id, err := uuid.Parse(activityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	activity, err := api.service.GetActivity(context.Background(), id)
	if err != nil || activity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	// Fetch Organization
	var organization *model.Organization
	if api.organizationRepo != nil {
		organization, _ = api.organizationRepo.GetByID(context.Background(), activity.OrgID)
	}

	// Fetch Badge
	var badge *model.Badge
	if activity.BadgeDefID != nil && api.badgeRepo != nil {
		badge, _ = api.badgeRepo.GetByID(context.Background(), *activity.BadgeDefID)
	}

	// Fetch Participations
	var participations []model.ActivityParticipation
	if api.participationService != nil {
		participations, _ = api.participationService.ListParticipations(context.Background(), &activity.ActivityID, nil, nil, 0, 100)
	}

	c.JSON(http.StatusOK, gin.H{
		"activity":       activity,
		"organization":   organization,
		"badge":          badge,
		"participations": participations,
	})
}

func (api *ActivityAPI) CreateActivity(c *gin.Context) {
	orgID := c.Param("org_id")
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activity := &model.Activity{
		ActivityID:   uuid.New(),
		OrgID:        orgUUID,
		ActivityName: req.ActivityName,
		Description:  &req.Description,
		Location:     &req.Location,
		// You may want to parse StartDate, EndDate, BadgeDefID as needed
	}
	err = api.service.CreateActivity(context.Background(), activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}
	c.JSON(http.StatusCreated, activity)
}

func (api *ActivityAPI) UpdateActivity(c *gin.Context) {
	activityID := c.Param("id")
	id, err := uuid.Parse(activityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{
		"activity_name": req.ActivityName,
		"description":   req.Description,
		"location":      req.Location,
		// Add other fields as needed
	}
	activity, err := api.service.UpdateActivity(context.Background(), id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func (api *ActivityAPI) DeleteActivity(c *gin.Context) {
	activityID := c.Param("id")
	id, err := uuid.Parse(activityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	err = api.service.DeleteActivity(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}
