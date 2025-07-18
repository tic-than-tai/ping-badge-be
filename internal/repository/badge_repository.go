package repository

import (
	"context"
	"ping-badge-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BadgeRepository interface {
	Create(ctx context.Context, badge *model.Badge) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Badge, error)
	List(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Badge, error)
	ListIssuedBadgesByUser(ctx context.Context, userID uuid.UUID) ([]model.IssuedBadge, error)
	CreateIssuedBadge(ctx context.Context, issuedBadge *model.IssuedBadge) error
	Update(ctx context.Context, badge *model.Badge) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type badgeRepositoryImpl struct {
	db *gorm.DB
}

func NewBadgeRepository(db *gorm.DB) BadgeRepository {
	return &badgeRepositoryImpl{db: db}
}

func (r *badgeRepositoryImpl) Create(ctx context.Context, badge *model.Badge) error {
	return r.db.WithContext(ctx).Create(badge).Error
}

func (r *badgeRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Badge, error) {
	var badge model.Badge
	err := r.db.WithContext(ctx).First(&badge, "badge_def_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &badge, nil
}

func (r *badgeRepositoryImpl) List(ctx context.Context, orgID *uuid.UUID, offset, limit int) ([]model.Badge, error) {
	var badges []model.Badge
	query := r.db.WithContext(ctx).Table("badges")
	if orgID != nil {
		query = query.Where("org_id = ?", *orgID)
	}
	err := query.Offset(offset).Limit(limit).Find(&badges).Error
	return badges, err
}

func (r *badgeRepositoryImpl) Update(ctx context.Context, badge *model.Badge) error {
	return r.db.WithContext(ctx).Table("badges").Save(badge).Error
}

func (r *badgeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Badge{}, "badge_def_id = ?", id).Error
}

func (r *badgeRepositoryImpl) ListIssuedBadgesByUser(ctx context.Context, userID uuid.UUID) ([]model.IssuedBadge, error) {
	var badges []model.IssuedBadge
	err := r.db.WithContext(ctx).Table("issued_badges").Where("user_id = ?", userID).Find(&badges).Error
	return badges, err
}

func (r *badgeRepositoryImpl) CreateIssuedBadge(ctx context.Context, issuedBadge *model.IssuedBadge) error {
	return r.db.WithContext(ctx).Create(issuedBadge).Error
}
