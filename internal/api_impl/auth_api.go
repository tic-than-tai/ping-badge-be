package api_impl

import (
	"context"
	"net/http"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthAPI struct {
	Service service.AuthService
}

func NewAuthAPI(service service.AuthService) *AuthAPI {
	return &AuthAPI{Service: service}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  model.User `json:"user"`
	Token string     `json:"token"`
}

func (api *AuthAPI) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := req.Role
	if role == "" {
		role = "USER"
	}
	user, token, err := api.Service.Register(context.Background(), req.Username, req.Email, req.Password, req.FullName, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, AuthResponse{User: *user, Token: token})
}

func (api *AuthAPI) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, token, err := api.Service.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, AuthResponse{User: *user, Token: token})
}

func (api *AuthAPI) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, err := api.Service.GetProfile(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (api *AuthAPI) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	var req struct {
		Username          string `json:"username"`
		FullName          string `json:"full_name"`
		ProfilePictureURL string `json:"profile_picture_url"`
		Bio               string `json:"bio"`
		PrivacySetting    string `json:"privacy_setting"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := api.Service.UpdateProfile(context.Background(), userID, req.Username, req.FullName, req.ProfilePictureURL, req.Bio, req.PrivacySetting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, user)
}
