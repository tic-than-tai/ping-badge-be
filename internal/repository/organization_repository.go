package repository

import (
	"context"
	"ping-badge-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	err := r.db.WithContext(ctx).First(&org, "org_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) List(ctx context.Context, offset, limit int, userID *uuid.UUID) ([]model.Organization, error) {
	var orgs []model.Organization
	query := r.db.WithContext(ctx)
	if userID != nil {
		query = query.Where("user_id_owner = ?", *userID)
	}
	err := query.Offset(offset).Limit(limit).Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) Update(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Organization{}, "org_id = ?", id).Error
}
