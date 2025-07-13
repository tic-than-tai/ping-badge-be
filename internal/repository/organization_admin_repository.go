package repository

import (
	"context"
	"ping-badge-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationAdminRepository interface {
	Create(ctx context.Context, admin *model.OrganizationAdmin) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.OrganizationAdmin, error)
	List(ctx context.Context, offset, limit int) ([]model.OrganizationAdmin, error)
	Update(ctx context.Context, admin *model.OrganizationAdmin) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type organizationAdminRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationAdminRepository(db *gorm.DB) OrganizationAdminRepository {
	return &organizationAdminRepositoryImpl{db: db}
}

func (r *organizationAdminRepositoryImpl) Create(ctx context.Context, admin *model.OrganizationAdmin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *organizationAdminRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.OrganizationAdmin, error) {
	var admin model.OrganizationAdmin
	err := r.db.WithContext(ctx).First(&admin, "admin_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *organizationAdminRepositoryImpl) List(ctx context.Context, offset, limit int) ([]model.OrganizationAdmin, error) {
	var admins []model.OrganizationAdmin
	err := r.db.WithContext(ctx).Table("organization_admins").Offset(offset).Limit(limit).Find(&admins).Error
	return admins, err
}

func (r *organizationAdminRepositoryImpl) Update(ctx context.Context, admin *model.OrganizationAdmin) error {
	return r.db.WithContext(ctx).Save(admin).Error
}

func (r *organizationAdminRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.OrganizationAdmin{}, "admin_id = ?", id).Error
}
