package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type BadgeService struct {
	repo *repository.BadgeRepository
}

func NewBadgeService(repo *repository.BadgeRepository) *BadgeService {
	return &BadgeService{repo: repo}
}

func (s *BadgeService) CreateBadge(ctx context.Context, badge *model.Badge) error {
	return s.repo.Create(ctx, badge)
}

func (s *BadgeService) GetBadge(ctx context.Context, id uuid.UUID) (*model.Badge, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BadgeService) ListBadges(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Badge, error) {
	return s.repo.List(ctx, orgID, offset, limit)
}

func (s *BadgeService) UpdateBadge(ctx context.Context, badge *model.Badge) error {
	return s.repo.Update(ctx, badge)
}

func (s *BadgeService) DeleteBadge(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
