package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type OrganizationService struct {
	repo *repository.OrganizationRepository
}

func NewOrganizationService(repo *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, org *model.Organization) error {
	return s.repo.Create(ctx, org)
}

func (s *OrganizationService) GetOrganization(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrganizationService) ListOrganizations(ctx context.Context, offset, limit int, userID *uuid.UUID) ([]model.Organization, error) {
	return s.repo.List(ctx, offset, limit, userID)
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, org *model.Organization) error {
	return s.repo.Update(ctx, org)
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
