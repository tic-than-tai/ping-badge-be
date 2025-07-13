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

type OrganizationAdminAPI struct {
	service service.OrganizationAdminService
}

func NewOrganizationAdminAPI(service service.OrganizationAdminService) *OrganizationAdminAPI {
	return &OrganizationAdminAPI{service: service}
}

type CreateOrganizationAdminRequest struct {
	OrgID  string `json:"org_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

func (api *OrganizationAdminAPI) CreateAdmin(c *gin.Context) {
	var req CreateOrganizationAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orgID, _ := uuid.Parse(req.OrgID)
	userID, _ := uuid.Parse(req.UserID)
	admin := &model.OrganizationAdmin{
		AdminID: uuid.New(),
		OrgID:   orgID,
		UserID:  userID,
		Role:    req.Role,
	}
	if err := api.service.CreateAdmin(context.Background(), admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}
	c.JSON(http.StatusCreated, admin)

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
	admins, err := api.service.ListAdmins(context.Background(), offset, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admins"})
		return
	}
	c.JSON(http.StatusOK, admins)
}

func (api *OrganizationAdminAPI) GetAdmin(c *gin.Context) {
	adminID := c.Param("id")
	id, err := uuid.Parse(adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}
	admin, err := api.service.GetAdmin(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
	c.JSON(http.StatusOK, admin)
}

func (api *OrganizationAdminAPI) UpdateAdmin(c *gin.Context) {
	adminID := c.Param("id")
	id, err := uuid.Parse(adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}
	var req CreateOrganizationAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admin, err := api.service.GetAdmin(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
	orgID, _ := uuid.Parse(req.OrgID)
	userID, _ := uuid.Parse(req.UserID)
	admin.OrgID = orgID
	admin.UserID = userID
	admin.Role = req.Role
	if err := api.service.UpdateAdmin(context.Background(), admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin"})
		return
	}
	c.JSON(http.StatusOK, admin)
}

func (api *OrganizationAdminAPI) DeleteAdmin(c *gin.Context) {
	adminID := c.Param("id")
	id, err := uuid.Parse(adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}
	if err := api.service.DeleteAdmin(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete admin"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})
}
