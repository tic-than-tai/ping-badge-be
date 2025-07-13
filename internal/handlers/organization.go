package handlers

import (
	"net/http"
	"strconv"

	"ping-badge-be/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationHandler struct {
	db *gorm.DB
}

func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{db: db}
}

type CreateOrgRequest struct {
	OrgName     string `json:"org_name" binding:"required"`
	OrgEmail    string `json:"org_email" binding:"required,email"`
	OrgLogoURL  string `json:"org_logo_url"`
	Description string `json:"description"`
	WebsiteURL  string `json:"website_url"`
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if organization name already exists
	var existingOrg models.Organization
	if err := h.db.Where("org_name = ? OR org_email = ?", req.OrgName, req.OrgEmail).First(&existingOrg).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Organization name or email already exists"})
		return
	}

	org := models.Organization{
		OrgID:       uuid.New(),
		OrgName:     req.OrgName,
		OrgEmail:    req.OrgEmail,
		UserIDOwner: userID.(uuid.UUID),
	}

	if req.OrgLogoURL != "" {
		org.OrgLogoURL = &req.OrgLogoURL
	}
	if req.Description != "" {
		org.Description = &req.Description
	}
	if req.WebsiteURL != "" {
		org.WebsiteURL = &req.WebsiteURL
	}

	if err := h.db.Create(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	// Load the organization with owner details
	if err := h.db.Preload("Owner").First(&org, org.OrgID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load organization details"})
		return
	}

	c.JSON(http.StatusCreated, org)
}

func (h *OrganizationHandler) GetOrganizations(c *gin.Context) {
	var organizations []models.Organization

	query := h.db.Preload("Owner")

	// Add pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	
	var pageInt, limitInt int = 1, 10
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageInt = p
	}
	if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
		limitInt = l
	}

	offset := (pageInt - 1) * limitInt

	if err := query.Offset(offset).Limit(limitInt).Find(&organizations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	c.JSON(http.StatusOK, organizations)
}

func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	orgID := c.Param("id")
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var org models.Organization
	if err := h.db.Preload("Owner").Preload("Admins.User").First(&org, "org_id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	orgID := c.Param("id")
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var org models.Organization
	if err := h.db.First(&org, "org_id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization"})
		return
	}

	// Check if user is owner
	if org.UserIDOwner != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner can update"})
		return
	}

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	updates := map[string]interface{}{
		"org_name":  req.OrgName,
		"org_email": req.OrgEmail,
	}

	if req.OrgLogoURL != "" {
		updates["org_logo_url"] = req.OrgLogoURL
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.WebsiteURL != "" {
		updates["website_url"] = req.WebsiteURL
	}

	if err := h.db.Model(&org).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	// Load updated organization
	if err := h.db.Preload("Owner").First(&org, "org_id = ?", orgID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	orgID := c.Param("id")
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var org models.Organization
	if err := h.db.First(&org, "org_id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization"})
		return
	}

	// Check if user is owner
	if org.UserIDOwner != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner can delete"})
		return
	}

	if err := h.db.Delete(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

func (h *OrganizationHandler) AddAdmin(c *gin.Context) {
	orgID := c.Param("id")
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if organization exists and user is owner
	var org models.Organization
	if err := h.db.First(&org, "org_id = ?", orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization"})
		return
	}

	if org.UserIDOwner != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only organization owner can add admins"})
		return
	}

	// Check if user exists
	var user models.User
	if err := h.db.First(&user, "user_id = ?", adminUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Check if admin already exists
	var existingAdmin models.OrganizationAdmin
	if err := h.db.Where("org_id = ? AND user_id = ?", orgID, adminUserID).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already an admin"})
		return
	}

	admin := models.OrganizationAdmin{
		AdminID: uuid.New(),
		OrgID:   uuid.MustParse(orgID),
		UserID:  adminUserID,
		Role:    req.Role,
	}

	if err := h.db.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add admin"})
		return
	}

	// Load admin with user details
	if err := h.db.Preload("User").First(&admin, admin.AdminID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load admin details"})
		return
	}

	c.JSON(http.StatusCreated, admin)
}
