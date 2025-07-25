package api_impl

import (
	"context"
	"net/http"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivityParticipationAPI struct {
	service service.ActivityParticipationService
}

func NewActivityParticipationAPI(service service.ActivityParticipationService) *ActivityParticipationAPI {
	return &ActivityParticipationAPI{service: service}
}

type CreateParticipationRequest struct {
	ActivityID              string  `json:"activity_id" binding:"required"`
	UserID                  string  `json:"user_id" binding:"required"`
	Status                  string  `json:"status"`
	ProofOfParticipationURL *string `json:"proof_of_participation_url"`
	IssuedBadgeID           *string `json:"issued_badge_id"`
}

func (api *ActivityParticipationAPI) ListParticipations(c *gin.Context) {
	activityID := c.Query("activity_id")
	userID := c.Query("user_id")
	// status := c.Query("status")

	var activityUUID, userUUID *uuid.UUID
	var statusPtr *string

	if activityID != "" {
		parsed, err := uuid.Parse(activityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
			return
		}
		activityUUID = &parsed
	}
	if userID != "" {
		parsed, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userUUID = &parsed
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
	participations, err := api.service.ListParticipations(context.Background(), activityUUID, userUUID, statusPtr, offset, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch participations"})
		return
	}
	c.JSON(http.StatusOK, participations)
}

func (api *ActivityParticipationAPI) GetParticipation(c *gin.Context) {
	id := c.Param("id")
	participationID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participation ID"})
		return
	}
	participation, err := api.service.GetParticipation(context.Background(), participationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Participation not found"})
		return
	}
	c.JSON(http.StatusOK, participation)
}

func (api *ActivityParticipationAPI) CreateParticipation(c *gin.Context) {
	var req CreateParticipationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activityID, err := uuid.Parse(req.ActivityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var issuedBadgeID *uuid.UUID
	if req.IssuedBadgeID != nil && *req.IssuedBadgeID != "" {
		parsed, err := uuid.Parse(*req.IssuedBadgeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issued badge ID"})
			return
		}
		issuedBadgeID = &parsed
	}
	participation := &model.ActivityParticipation{
		ActivityID:              activityID,
		UserID:                  userID,
		Status:                  req.Status,
		ProofOfParticipationURL: req.ProofOfParticipationURL,
		IssuedBadgeID:           issuedBadgeID,
		CreatedAt:               time.Now(),
	}
	if err := api.service.CreateParticipation(context.Background(), participation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create participation"})
		return
	}
	c.JSON(http.StatusCreated, participation)
}

func (api *ActivityParticipationAPI) UpdateParticipation(c *gin.Context) {
	id := c.Param("id")
	participationID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participation ID"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	participation, err := api.service.UpdateParticipation(context.Background(), participationID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update participation"})
		return
	}
	c.JSON(http.StatusOK, participation)
}

func (api *ActivityParticipationAPI) DeleteParticipation(c *gin.Context) {
	id := c.Param("id")
	participationID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participation ID"})
		return
	}
	if err := api.service.DeleteParticipation(context.Background(), participationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete participation"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

type UploadEvidenceRequest struct {
	ProofOfParticipationURL string `json:"proof_of_participation_url" binding:"required"`
}

func (api *ActivityParticipationAPI) UploadEvidence(c *gin.Context) {
	participationIDStr := c.Param("id")
	userIDStr := c.Query("user_id")
	activityIDStr := c.Query("activity_id")

	var participationID *uuid.UUID
	if participationIDStr != "" {
		parsed, err := uuid.Parse(participationIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participation ID"})
			return
		}
		participationID = &parsed
	}
	var userID *uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userID = &parsed
	}
	var activityID *uuid.UUID
	if activityIDStr != "" {
		parsed, err := uuid.Parse(activityIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
			return
		}
		activityID = &parsed
	}

	var req UploadEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find participation by filters
	var participation *model.ActivityParticipation
	var err error
	if participationID != nil {
		participation, err = api.service.GetParticipation(context.Background(), *participationID)
	} else {
		var status = "NOT_COMPLETED"
		participations, err := api.service.ListParticipations(context.Background(), activityID, userID, &status, 1, 1)
		if err != nil || len(participations) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Participation not found"})
			return
		}
		participation = &participations[0]
	}
	if err != nil || participation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Participation not found"})
		return
	}

	updates := map[string]interface{}{
		"proof_of_participation_url": req.ProofOfParticipationURL,
	}
	updated, err := api.service.UpdateParticipation(context.Background(), participation.ParticipationID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload evidence"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

type UpdateParticipationStatusRequest struct {
	ProofOfParticipationURL *string `json:"proof_of_participation_url"`
	Status                  string  `json:"status" binding:"required"`
}

func (api *ActivityParticipationAPI) UpdateParticipationStatus(c *gin.Context) {
	id := c.Param("id")
	participationID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participation ID"})
		return
	}

	var req UpdateParticipationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participation, err := api.service.UpdateParticipationWithBadgeCreation(
		context.Background(),
		participationID,
		req.ProofOfParticipationURL,
		strings.ToUpper(req.Status),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update participation status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Participation updated successfully",
		"participation": participation,
		"badge_created": req.Status == "completed" || req.Status == "COMPLETED",
	})
}
