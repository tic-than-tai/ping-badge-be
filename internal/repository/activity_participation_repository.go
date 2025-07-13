package repository

import (
	"ping-badge-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityParticipationRepository interface {
	Create(participation *model.ActivityParticipation) error
	FindByID(id uuid.UUID) (*model.ActivityParticipation, error)
	FindAll(activityID *uuid.UUID, userID *uuid.UUID, offset, limit int) ([]model.ActivityParticipation, error)
	Update(id uuid.UUID, updates map[string]interface{}) (*model.ActivityParticipation, error)
	Delete(id uuid.UUID) error
}

type activityParticipationRepositoryImpl struct {
	db *gorm.DB
}

func NewActivityParticipationRepository(db *gorm.DB) ActivityParticipationRepository {
	return &activityParticipationRepositoryImpl{db: db}
}

func (r *activityParticipationRepositoryImpl) Create(participation *model.ActivityParticipation) error {
	return r.db.Create(participation).Error
}

func (r *activityParticipationRepositoryImpl) FindByID(id uuid.UUID) (*model.ActivityParticipation, error) {
	var participation model.ActivityParticipation
	err := r.db.First(&participation, "participation_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &participation, nil
}

func (r *activityParticipationRepositoryImpl) FindAll(activityID *uuid.UUID, userID *uuid.UUID, offset, limit int) ([]model.ActivityParticipation, error) {
	var participations []model.ActivityParticipation
	query := r.db
	if activityID != nil {
		query = query.Where("activity_id = ?", *activityID)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.Offset(offset).Limit(limit).Find(&participations).Error
	return participations, err
}

func (r *activityParticipationRepositoryImpl) Update(id uuid.UUID, updates map[string]interface{}) (*model.ActivityParticipation, error) {
	var participation model.ActivityParticipation
	err := r.db.First(&participation, "participation_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&participation).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &participation, nil
}

func (r *activityParticipationRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ActivityParticipation{}, "participation_id = ?", id).Error
}
