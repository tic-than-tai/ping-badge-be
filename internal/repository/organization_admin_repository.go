package repository

import (
	"context"
	"ping-badge-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationAdminRepository struct {
	db *gorm.DB
}

func NewOrganizationAdminRepository(db *gorm.DB) *OrganizationAdminRepository {
	return &OrganizationAdminRepository{db: db}
}

func (r *OrganizationAdminRepository) Create(ctx context.Context, admin *model.OrganizationAdmin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *OrganizationAdminRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.OrganizationAdmin, error) {
	var admin model.OrganizationAdmin
	err := r.db.WithContext(ctx).First(&admin, "admin_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *OrganizationAdminRepository) List(ctx context.Context, offset, limit int) ([]model.OrganizationAdmin, error) {
	var admins []model.OrganizationAdmin
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&admins).Error
	return admins, err
}

func (r *OrganizationAdminRepository) Update(ctx context.Context, admin *model.OrganizationAdmin) error {
	return r.db.WithContext(ctx).Save(admin).Error
}

func (r *OrganizationAdminRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.OrganizationAdmin{}, "admin_id = ?", id).Error
}
