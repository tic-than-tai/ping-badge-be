package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type BadgeService interface {
	CreateBadge(ctx context.Context, badge *model.Badge) error
	GetBadge(ctx context.Context, id uuid.UUID) (*model.Badge, error)
	ListBadges(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Badge, error)
	UpdateBadge(ctx context.Context, badge *model.Badge) error
	DeleteBadge(ctx context.Context, id uuid.UUID) error
}

type badgeServiceImpl struct {
	repo repository.BadgeRepository
}

func NewBadgeService(repo repository.BadgeRepository) BadgeService {
	return &badgeServiceImpl{repo: repo}
}

func (s *badgeServiceImpl) CreateBadge(ctx context.Context, badge *model.Badge) error {
	// Add business logic, validation, authorization here
	return s.repo.Create(ctx, badge)
}

func (s *badgeServiceImpl) GetBadge(ctx context.Context, id uuid.UUID) (*model.Badge, error) {
	// Add business logic, validation, authorization here
	return s.repo.GetByID(ctx, id)
}

func (s *badgeServiceImpl) ListBadges(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Badge, error) {
	// Add business logic, validation, authorization here
	return s.repo.List(ctx, orgID, offset, limit)
}

func (s *badgeServiceImpl) UpdateBadge(ctx context.Context, badge *model.Badge) error {
	// Add business logic, validation, authorization here
	return s.repo.Update(ctx, badge)
}

func (s *badgeServiceImpl) DeleteBadge(ctx context.Context, id uuid.UUID) error {
	// Add business logic, validation, authorization here
	return s.repo.Delete(ctx, id)
}
