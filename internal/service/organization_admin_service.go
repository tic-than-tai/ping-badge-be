package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type OrganizationAdminService interface {
	CreateAdmin(ctx context.Context, admin *model.OrganizationAdmin) error
	GetAdmin(ctx context.Context, id uuid.UUID) (*model.OrganizationAdmin, error)
	ListAdmins(ctx context.Context, offset, limit int) ([]model.OrganizationAdmin, error)
	UpdateAdmin(ctx context.Context, admin *model.OrganizationAdmin) error
	DeleteAdmin(ctx context.Context, id uuid.UUID) error
}

type organizationAdminServiceImpl struct {
	repo repository.OrganizationAdminRepository
}

func NewOrganizationAdminService(repo repository.OrganizationAdminRepository) OrganizationAdminService {
	return &organizationAdminServiceImpl{repo: repo}
}

func (s *organizationAdminServiceImpl) CreateAdmin(ctx context.Context, admin *model.OrganizationAdmin) error {
	// Add business logic, validation, authorization here
	return s.repo.Create(ctx, admin)
}

func (s *organizationAdminServiceImpl) GetAdmin(ctx context.Context, id uuid.UUID) (*model.OrganizationAdmin, error) {
	// Add business logic, validation, authorization here
	return s.repo.GetByID(ctx, id)
}

func (s *organizationAdminServiceImpl) ListAdmins(ctx context.Context, offset, limit int) ([]model.OrganizationAdmin, error) {
	// Add business logic, validation, authorization here
	return s.repo.List(ctx, offset, limit)
}

func (s *organizationAdminServiceImpl) UpdateAdmin(ctx context.Context, admin *model.OrganizationAdmin) error {
	// Add business logic, validation, authorization here
	return s.repo.Update(ctx, admin)
}

func (s *organizationAdminServiceImpl) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	// Add business logic, validation, authorization here
	return s.repo.Delete(ctx, id)
}
