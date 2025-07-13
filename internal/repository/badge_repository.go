package repository

import (
	"context"
	"ping-badge-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BadgeRepository struct {
	db *gorm.DB
}

func NewBadgeRepository(db *gorm.DB) *BadgeRepository {
	return &BadgeRepository{db: db}
}

func (r *BadgeRepository) Create(ctx context.Context, badge *model.Badge) error {
	return r.db.WithContext(ctx).Create(badge).Error
}

func (r *BadgeRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Badge, error) {
	var badge model.Badge
	err := r.db.WithContext(ctx).First(&badge, "badge_def_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &badge, nil
}

func (r *BadgeRepository) List(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Badge, error) {
	var badges []model.Badge
	query := r.db.WithContext(ctx)
	if orgID != nil {
		query = query.Where("org_id = ?", *orgID)
	}
	err := query.Offset(offset).Limit(limit).Find(&badges).Error
	return badges, err
}

func (r *BadgeRepository) Update(ctx context.Context, badge *model.Badge) error {
	return r.db.WithContext(ctx).Save(badge).Error
}

func (r *BadgeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Badge{}, "badge_def_id = ?", id).Error
}
