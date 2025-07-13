package api_impl

import (
	"context"
	"net/http"

	"ping-badge-be/internal/model"
	"ping-badge-be/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationAPI struct {
	service *service.OrganizationService
}

func NewOrganizationAPI(service *service.OrganizationService) *OrganizationAPI {
	return &OrganizationAPI{service: service}
}

type CreateOrganizationRequest struct {
	OrgName     string `json:"org_name" binding:"required"`
	OrgEmail    string `json:"org_email" binding:"required"`
	Description string `json:"description"`
	WebsiteURL  string `json:"website_url"`
}

func (api *OrganizationAPI) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	org := &model.Organization{
		OrgID:       uuid.New(),
		OrgName:     req.OrgName,
		OrgEmail:    req.OrgEmail,
		Description: &req.Description,
		WebsiteURL:  &req.WebsiteURL,
	}
	if err := api.service.CreateOrganization(context.Background(), org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}
	c.JSON(http.StatusCreated, org)
}

// Alias for router: ListOrganizations
func (api *OrganizationAPI) ListOrganizations(c *gin.Context) {
	// Call the service to list organizations
	offset := 0
	limit := 10
	orgs, err := api.service.ListOrganizations(context.Background(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

// AddAdmin endpoint stub for router
func (api *OrganizationAPI) AddAdmin(c *gin.Context) {
	// Validate request, call service, return response
	c.JSON(http.StatusNotImplemented, gin.H{"error": "AddAdmin not implemented yet"})
}

func (api *OrganizationAPI) GetOrganization(c *gin.Context) {
	orgID := c.Param("id")
	id, err := uuid.Parse(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	org, err := api.service.GetOrganization(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (api *OrganizationAPI) UpdateOrganization(c *gin.Context) {
	orgID := c.Param("id")
	id, err := uuid.Parse(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	org, err := api.service.GetOrganization(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	org.OrgName = req.OrgName
	org.OrgEmail = req.OrgEmail
	org.Description = &req.Description
	org.WebsiteURL = &req.WebsiteURL
	if err := api.service.UpdateOrganization(context.Background(), org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (api *OrganizationAPI) DeleteOrganization(c *gin.Context) {
	orgID := c.Param("id")
	id, err := uuid.Parse(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	if err := api.service.DeleteOrganization(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}
