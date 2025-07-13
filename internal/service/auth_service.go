package service

import (
	"context"
	"ping-badge-be/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password, fullName, role string) (*model.User, string, error)
	Login(ctx context.Context, email, password string) (*model.User, string, error)
	GetProfile(ctx context.Context, userID interface{}) (*model.User, error)
	UpdateProfile(ctx context.Context, userID interface{}, username, fullName, profilePictureURL, bio, privacySetting string) (*model.User, error)
}
