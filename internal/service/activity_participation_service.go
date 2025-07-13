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
	ListParticipations(ctx context.Context, activityID *uuid.UUID, userID *uuid.UUID, offset, limit int) ([]model.ActivityParticipation, error)
	UpdateParticipation(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*model.ActivityParticipation, error)
	DeleteParticipation(ctx context.Context, id uuid.UUID) error
}

type activityParticipationServiceImpl struct {
	repo repository.ActivityParticipationRepository
}

func NewActivityParticipationService(repo repository.ActivityParticipationRepository) ActivityParticipationService {
	return &activityParticipationServiceImpl{repo: repo}
}

func (s *activityParticipationServiceImpl) CreateParticipation(ctx context.Context, participation *model.ActivityParticipation) error {
	// Add business logic, validation, authorization here
	return s.repo.Create(participation)
}

func (s *activityParticipationServiceImpl) GetParticipation(ctx context.Context, id uuid.UUID) (*model.ActivityParticipation, error) {
	return s.repo.FindByID(id)
}

func (s *activityParticipationServiceImpl) ListParticipations(ctx context.Context, activityID *uuid.UUID, userID *uuid.UUID, offset, limit int) ([]model.ActivityParticipation, error) {
	return s.repo.FindAll(activityID, userID, offset, limit)
}

func (s *activityParticipationServiceImpl) UpdateParticipation(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*model.ActivityParticipation, error) {
	return s.repo.Update(id, updates)
}

func (s *activityParticipationServiceImpl) DeleteParticipation(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}
