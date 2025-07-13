package repository

import (
	"ping-badge-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(activity *model.Activity) error
	FindByID(activityID uuid.UUID) (*model.Activity, error)
	FindAll(orgID *uuid.UUID, offset, limit int) ([]model.Activity, error)
	Update(activityID uuid.UUID, updates map[string]interface{}) (*model.Activity, error)
	Delete(activityID uuid.UUID) error
}

// ActivityRepositoryImpl implements ActivityRepository
// Implementation will be added after interface usage
type activityRepositoryImpl struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepositoryImpl{db: db}
}

func (r *activityRepositoryImpl) Create(activity *model.Activity) error {
	return r.db.Create(activity).Error
}

func (r *activityRepositoryImpl) FindByID(activityID uuid.UUID) (*model.Activity, error) {
	var activity model.Activity
	err := r.db.First(&activity, "activity_id = ?", activityID).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepositoryImpl) FindAll(orgID *uuid.UUID, offset, limit int) ([]model.Activity, error) {
	var activities []model.Activity
	query := r.db
	if orgID != nil {
		query = query.Where("org_id = ?", *orgID)
	}
	err := query.Offset(offset).Limit(limit).Find(&activities).Error
	return activities, err
}

func (r *activityRepositoryImpl) Update(activityID uuid.UUID, updates map[string]interface{}) (*model.Activity, error) {
	var activity model.Activity
	err := r.db.First(&activity, "activity_id = ?", activityID).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&activity).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepositoryImpl) Delete(activityID uuid.UUID) error {
	return r.db.Delete(&model.Activity{}, "activity_id = ?", activityID).Error
}

// Implement participation methods as needed
