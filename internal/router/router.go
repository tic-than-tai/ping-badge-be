package router

import (
	"ping-badge-be/internal/api_impl"
	"ping-badge-be/internal/config"
	"ping-badge-be/internal/middleware"
	"ping-badge-be/internal/repository"
	"ping-badge-be/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS(cfg.CORSOrigins))

	// Initialize Auth API (layered architecture)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authAPI := api_impl.NewAuthAPI(authService)

	// Initialize User API (layered architecture)
	userService := service.NewUserService(userRepo)
	userAPI := api_impl.NewUserAPI(userService)

	// Initialize OrganizationAdmin API (layered architecture)
	orgAdminRepo := repository.NewOrganizationAdminRepository(db)
	orgAdminService := service.NewOrganizationAdminService(orgAdminRepo)
	orgAdminAPI := api_impl.NewOrganizationAdminAPI(orgAdminService)
	// orgHandler removed: use OrganizationAPI for all organization routes (layered architecture)
	orgRepo := repository.NewOrganizationRepository(db)
	orgService := service.NewOrganizationService(orgRepo)
	orgAPI := api_impl.NewOrganizationAPI(orgService)

	// Initialize Badge API (layered architecture)
	badgeRepo := repository.NewBadgeRepository(db)
	badgeService := service.NewBadgeService(badgeRepo)
	badgeAPI := api_impl.NewBadgeAPI(badgeService)

	// Initialize Activity API (layered architecture)
	activityRepo := repository.NewActivityRepository(db)
	activityService := service.NewActivityService(activityRepo)
	activityAPI := api_impl.NewActivityAPI(activityService)

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
			auth.POST("/register", authAPI.Register)
			auth.POST("/login", authAPI.Login)
		}

		// Public organization routes (use OrganizationAPI)
		api.GET("/organizations", orgAPI.ListOrganizations)
		api.GET("/organizations/:id", orgAPI.GetOrganization)

		// Public badge routes (use BadgeAPI)
		api.GET("/badges", badgeAPI.ListBadges)
		api.GET("/badges/:id", badgeAPI.GetBadge)

		// Public activity routes (use ActivityAPI)
		api.GET("/activities", activityAPI.ListActivities)
		api.GET("/activities/:id", activityAPI.GetActivity)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// User routes (use UserAPI)
		protected.GET("/users", userAPI.ListUsers)
		protected.GET("/users/:id", userAPI.GetUser)
		protected.POST("/users", userAPI.CreateUser)
		protected.PUT("/users/:id", userAPI.UpdateUser)
		protected.DELETE("/users/:id", userAPI.DeleteUser)

		// OrganizationAdmin routes (use OrganizationAdminAPI)
		protected.POST("/organizations/:id/admins", orgAdminAPI.CreateAdmin)
		protected.GET("/organizations/:id/admins", orgAdminAPI.GetAdmin)
		protected.PUT("/organizations/:id/admins/:admin_id", orgAdminAPI.UpdateAdmin)
		protected.DELETE("/organizations/:id/admins/:admin_id", orgAdminAPI.DeleteAdmin)

		// Badge routes (use BadgeAPI)
		protected.POST("/organizations/:id/badges", badgeAPI.CreateBadge)
		protected.PUT("/badges/:id", badgeAPI.UpdateBadge)
		protected.DELETE("/badges/:id", badgeAPI.DeleteBadge)
		// protected.POST("/badges/:id/issue", badgeAPI.IssueBadge) // Add to BadgeAPI if needed

		// Activity routes (use ActivityAPI)
		protected.POST("/organizations/:id/activities", activityAPI.CreateActivity)
		protected.PUT("/activities/:id", activityAPI.UpdateActivity)
		protected.DELETE("/activities/:id", activityAPI.DeleteActivity)
		// Add join and participations endpoints to ActivityAPI as needed

		// Auth protected routes
		protected.GET("/auth/profile", authAPI.GetProfile)
		protected.PUT("/auth/profile", authAPI.UpdateProfile)
	}

	return r
}
