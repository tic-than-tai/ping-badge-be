package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type ActivityService interface {
	CreateActivity(ctx context.Context, activity *model.Activity) error
	GetActivity(ctx context.Context, id uuid.UUID) (*model.Activity, error)
	ListActivities(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Activity, error)
	UpdateActivity(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*model.Activity, error)
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type activityServiceImpl struct {
	repo repository.ActivityRepository
}

func NewActivityService(repo repository.ActivityRepository) ActivityService {
	return &activityServiceImpl{repo: repo}
}

func (s *activityServiceImpl) CreateActivity(ctx context.Context, activity *model.Activity) error {
	// Add business logic, validation, authorization here
	return s.repo.Create(activity)
}

func (s *activityServiceImpl) GetActivity(ctx context.Context, id uuid.UUID) (*model.Activity, error) {
	// Add business logic, validation, authorization here
	return s.repo.FindByID(id)
}

func (s *activityServiceImpl) ListActivities(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Activity, error) {
	// Add business logic, validation, authorization here
	return s.repo.FindAll(orgID, offset, limit)
}

func (s *activityServiceImpl) UpdateActivity(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*model.Activity, error) {
	// Add business logic, validation, authorization here
	return s.repo.Update(id, updates)
}

func (s *activityServiceImpl) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	// Add business logic, validation, authorization here
	return s.repo.Delete(id)
}
