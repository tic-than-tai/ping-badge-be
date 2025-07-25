package api_impl

import (
	"context"
	"net/http"
	"ping-badge-be/internal/constant"
	"ping-badge-be/internal/model"
	"ping-badge-be/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserAPI struct {
	service service.UserService
}

func NewUserAPI(service service.UserService) *UserAPI {
	return &UserAPI{service: service}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	FullName string `json:"full_name"`
}

func (api *UserAPI) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &model.User{
		UserID:       uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password, // Hash in real impl
	}
	if req.FullName != "" {
		user.FullName = &req.FullName
	}
	if err := api.service.CreateUser(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (api *UserAPI) GetUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := api.service.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (api *UserAPI) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	page := c.DefaultQuery("page", strconv.Itoa(constant.DefaultPage))
	limit := c.DefaultQuery("limit", strconv.Itoa(constant.DefaultLimit))
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	if pageInt < 1 {
		pageInt = constant.DefaultPage
	}
	if limitInt < 1 || limitInt > constant.MaxLimit {
		limitInt = constant.DefaultLimit
	}
	offset := (pageInt - 1) * limitInt
	users, err := api.service.ListUsers(context.Background(), offset, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (api *UserAPI) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := api.service.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Username = req.Username
	user.Email = req.Email
	if req.FullName != "" {
		user.FullName = &req.FullName
	}
	if err := api.service.UpdateUser(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (api *UserAPI) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := api.service.DeleteUser(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
