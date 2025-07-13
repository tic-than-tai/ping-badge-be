package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type ActivityParticipationService interface {
	CreateParticipation(ctx context.Context, participation *model.ActivityParticipation) error
	GetParticipation(ctx context.Context, id uuid.UUID) (*model.ActivityParticipation, error)
	ListParticipations(ctx context.Context, activityID *uuid.UUID, userID *uuid.UUID, status *string, offset, limit int) ([]model.ActivityParticipation, error)
	UpdateParticipation(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*model.ActivityParticipation, error)
	UpdateParticipationWithBadgeCreation(ctx context.Context, id uuid.UUID, proofURL *string, status string) (*model.ActivityParticipation, error)
	DeleteParticipation(ctx context.Context, id uuid.UUID) error
}

type activityParticipationServiceImpl struct {
	repo         repository.ActivityParticipationRepository
	activityRepo repository.ActivityRepository
	badgeRepo    repository.BadgeRepository
}

func NewActivityParticipationService(repo repository.ActivityParticipationRepository, activityRepo repository.ActivityRepository, badgeRepo repository.BadgeRepository) ActivityParticipationService {
	return &activityParticipationServiceImpl{
		repo:         repo,
		activityRepo: activityRepo,
		badgeRepo:    badgeRepo,
	}
}

func (s *activityParticipationServiceImpl) CreateParticipation(ctx context.Context, participation *model.ActivityParticipation) error {
	// Add business logic, validation, authorization here
	return s.repo.Create(participation)
}

func (s *activityParticipationServiceImpl) GetParticipation(ctx context.Context, id uuid.UUID) (*model.ActivityParticipation, error) {
	return s.repo.FindByID(id)
}

func (s *activityParticipationServiceImpl) ListParticipations(ctx context.Context, activityID *uuid.UUID, userID *uuid.UUID, status *string, offset, limit int) ([]model.ActivityParticipation, error) {
	return s.repo.FindAll(activityID, userID, status, offset, limit)
}

func (s *activityParticipationServiceImpl) UpdateParticipation(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*model.ActivityParticipation, error) {
	return s.repo.Update(id, updates)
}

func (s *activityParticipationServiceImpl) DeleteParticipation(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *activityParticipationServiceImpl) UpdateParticipationWithBadgeCreation(ctx context.Context, id uuid.UUID, proofURL *string, status string) (*model.ActivityParticipation, error) {
	// Get the current participation
	_, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if proofURL != nil {
		updates["proof_of_participation_url"] = *proofURL
	}
	if status != "" {
		updates["status"] = status
	}

	// Update the participation
	updatedParticipation, err := s.repo.Update(id, updates)
	if err != nil {
		return nil, err
	}

	// If status is COMPLETED, create a badge
	if status == "COMPLETED" {
		err := s.createBadgeForCompletion(ctx, updatedParticipation)
		if err != nil {
			// Log error but don't fail the update
			// In production, you might want to use proper logging
			// log.Printf("Failed to create badge for participation %s: %v", id, err)
		}
	}

	return updatedParticipation, nil
}

func (s *activityParticipationServiceImpl) createBadgeForCompletion(ctx context.Context, participation *model.ActivityParticipation) error {
	// Get the activity to find the associated badge
	activity, err := s.activityRepo.FindByID(participation.ActivityID)
	if err != nil {
		return err
	}

	// Check if activity has an associated badge
	if activity.BadgeDefID == nil {
		// No badge associated with this activity
		return nil
	}

	// Check if badge already exists for this user and activity
	existingBadges, err := s.badgeRepo.ListIssuedBadgesByUser(ctx, participation.UserID)
	if err != nil {
		return err
	}

	// Check if user already has this badge
	for _, badge := range existingBadges {
		if badge.BadgeDefID == *activity.BadgeDefID && badge.SourceID != nil && *badge.SourceID == participation.ActivityID {
			// Badge already exists for this activity
			return nil
		}
	}

	// Create new issued badge
	issuedBadge := &model.IssuedBadge{
		IssuedBadgeID:    uuid.New(),
		BadgeDefID:       *activity.BadgeDefID,
		UserID:           participation.UserID,
		OrgID:            activity.OrgID,
		VerificationCode: generateVerificationCode(),
		SourceType:       stringPtr("activity"),
		SourceID:         &participation.ActivityID,
		Status:           "issued",
	}

	// Create the issued badge
	err = s.badgeRepo.CreateIssuedBadge(ctx, issuedBadge)
	if err != nil {
		return err
	}

	// Update participation with issued badge ID
	_, err = s.repo.Update(participation.ParticipationID, map[string]interface{}{
		"issued_badge_id": issuedBadge.IssuedBadgeID,
	})

	return err
}

func generateVerificationCode() string {
	// Simple verification code generation
	// In production, you might want a more sophisticated approach
	return "VERIFY-" + uuid.New().String()[:8]
}

func stringPtr(s string) *string {
	return &s
}
