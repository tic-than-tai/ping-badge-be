package api_impl

import (
	"context"
	"net/http"
	"ping-badge-be/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserStatisticsAPI struct {
	badgeService         service.BadgeService
	activityService      service.ActivityService
	participationService service.ActivityParticipationService
}

func NewUserStatisticsAPI(badgeService service.BadgeService, activityService service.ActivityService, participationService service.ActivityParticipationService) *UserStatisticsAPI {
	return &UserStatisticsAPI{
		badgeService:         badgeService,
		activityService:      activityService,
		participationService: participationService,
	}
}

type UserStatisticsResponse struct {
	TotalBadges         int `json:"total_badges"`
	ActivitiesCompleted int `json:"activities_completed"`
	UpcomingEvents      int `json:"upcoming_events"`
}

func (api *UserStatisticsAPI) GetUserStatistics(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Total Badges
	badges, err := api.badgeService.ListIssuedBadgesByUser(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badges"})
		return
	}
	totalBadges := len(badges)

	// Activities Completed
	participations, err := api.participationService.ListParticipations(context.Background(), nil, &userID, 0, 1000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch participations"})
		return
	}
	activitiesCompleted := 0
	for _, p := range participations {
		activity, err := api.activityService.GetActivity(context.Background(), p.ActivityID)
		if err != nil || activity.StartDate == nil {
			continue
		}
		if p.CreatedAt.Before(*activity.StartDate) {
			activitiesCompleted++
		}
	}

	// Upcoming Events
	upcomingEvents := 0
	now := time.Now()
	for _, p := range participations {
		activity, err := api.activityService.GetActivity(context.Background(), p.ActivityID)
		if err != nil || activity.StartDate == nil || activity.EndDate == nil {
			continue
		}
		if now.Before(*activity.EndDate) && now.After(*activity.StartDate) {
			upcomingEvents++
		}
	}

	resp := UserStatisticsResponse{
		TotalBadges:         totalBadges,
		ActivitiesCompleted: activitiesCompleted,
		UpcomingEvents:      upcomingEvents,
	}
	c.JSON(http.StatusOK, resp)
}
