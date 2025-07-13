package service

import (
	"context"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	ListUsers(ctx context.Context, offset, limit int) ([]model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *model.User) error {
	// Add business logic, validation, authorization here
	return s.repo.Create(ctx, user)
}

func (s *userServiceImpl) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	// Add business logic, validation, authorization here
	return s.repo.GetByID(ctx, id)
}

func (s *userServiceImpl) ListUsers(ctx context.Context, offset, limit int) ([]model.User, error) {
	// Add business logic, validation, authorization here
	return s.repo.List(ctx, offset, limit)
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, user *model.User) error {
	// Add business logic, validation, authorization here
	return s.repo.Update(ctx, user)
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Add business logic, validation, authorization here
	return s.repo.Delete(ctx, id)
}
