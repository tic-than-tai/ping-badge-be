package router

import (
	"ping-badge-be/internal/config"
	"ping-badge-be/internal/handlers"
	"ping-badge-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS(cfg.CORSOrigins))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	orgHandler := handlers.NewOrganizationHandler(db)
	badgeHandler := handlers.NewBadgeHandler(db)
	activityHandler := handlers.NewActivityHandler(db)

	// Public routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public organization routes
		api.GET("/organizations", orgHandler.GetOrganizations)
		api.GET("/organizations/:id", orgHandler.GetOrganization)

		// Public badge routes
		api.GET("/badges", badgeHandler.GetBadges)
		api.GET("/badges/:id", badgeHandler.GetBadge)

		// Public activity routes
		api.GET("/activities", activityHandler.GetActivities)
		api.GET("/activities/:id", activityHandler.GetActivity)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// User profile routes
		protected.GET("/profile", authHandler.GetProfile)
		protected.PUT("/profile", authHandler.UpdateProfile)

		// Organization routes
		protected.POST("/organizations", orgHandler.CreateOrganization)
		protected.PUT("/organizations/:id", orgHandler.UpdateOrganization)
		protected.DELETE("/organizations/:id", orgHandler.DeleteOrganization)
		protected.POST("/organizations/:id/admins", orgHandler.AddAdmin)

		// Badge routes
		protected.POST("/organizations/:org_id/badges", badgeHandler.CreateBadge)
		protected.PUT("/badges/:id", badgeHandler.UpdateBadge)
		protected.DELETE("/badges/:id", badgeHandler.DeleteBadge)
		protected.POST("/badges/:id/issue", badgeHandler.IssueBadge)

		// Activity routes
		protected.POST("/organizations/:org_id/activities", activityHandler.CreateActivity)
		protected.PUT("/activities/:id", activityHandler.UpdateActivity)
		protected.DELETE("/activities/:id", activityHandler.DeleteActivity)
		protected.POST("/activities/:id/join", activityHandler.JoinActivity)
		protected.GET("/my-participations", activityHandler.GetUserParticipations)
	}

	return r
}
