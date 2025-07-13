package service

import (
	"context"
	"errors"
	"ping-badge-be/internal/middleware"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) AuthService {
	return &AuthServiceImpl{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthServiceImpl) Register(ctx context.Context, username, email, password, fullName, role string) (*model.User, string, error) {
	// Check if user already exists
	existing, err := s.repo.FindByEmailOrUsername(ctx, email, username)
	if err == nil && existing != nil {
		return nil, "", ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	if role == "" {
		role = "USER"
	}

	user := &model.User{
		UserID:       uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}
	if fullName != "" {
		user.FullName = &fullName
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}

	token, err := middleware.GenerateToken(user.UserID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}
	user.PasswordHash = ""
	return user, token, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}
	token, err := middleware.GenerateToken(user.UserID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}
	user.PasswordHash = ""
	return user, token, nil
}

func (s *AuthServiceImpl) GetProfile(ctx context.Context, userID interface{}) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthServiceImpl) UpdateProfile(ctx context.Context, userID interface{}, username, fullName, profilePictureURL, bio, privacySetting string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	if username != "" {
		user.Username = username
	}
	if fullName != "" {
		user.FullName = &fullName
	}
	if profilePictureURL != "" {
		user.ProfilePictureURL = &profilePictureURL
	}
	if bio != "" {
		user.Bio = &bio
	}
	if privacySetting != "" {
		user.PrivacySetting = privacySetting
	}
	user.UpdatedAt = time.Now()
	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

// Error definitions
var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)
