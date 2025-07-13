package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type OrganizationAdminService struct {
	repo *repository.OrganizationAdminRepository
}

func NewOrganizationAdminService(repo *repository.OrganizationAdminRepository) *OrganizationAdminService {
	return &OrganizationAdminService{repo: repo}
}

func (s *OrganizationAdminService) CreateAdmin(ctx context.Context, admin *model.OrganizationAdmin) error {
	return s.repo.Create(ctx, admin)
}

func (s *OrganizationAdminService) GetAdmin(ctx context.Context, id uuid.UUID) (*model.OrganizationAdmin, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrganizationAdminService) ListAdmins(ctx context.Context, offset, limit int) ([]model.OrganizationAdmin, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *OrganizationAdminService) UpdateAdmin(ctx context.Context, admin *model.OrganizationAdmin) error {
	return s.repo.Update(ctx, admin)
}

func (s *OrganizationAdminService) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
