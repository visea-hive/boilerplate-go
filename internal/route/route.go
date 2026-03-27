package route

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/visea-hive/auth-core/internal/handler"
	"github.com/visea-hive/auth-core/internal/middleware"
	"github.com/visea-hive/auth-core/internal/repository"
	"github.com/visea-hive/auth-core/internal/service"
	"github.com/visea-hive/auth-core/pkg/crypto"
	jwt "github.com/visea-hive/auth-core/pkg/jwt"
	"gorm.io/gorm"
)

func InitRoutes(router *gin.Engine, rdb *redis.Client, db *gorm.DB, jwtManager *jwt.Manager, hasher *crypto.Hasher) {
	// Initialize Dependencies
	orgRepo := repository.NewOrganizationRepository(db)
	orgService := service.NewOrganizationService(orgRepo, db)
	orgHandler := handler.NewOrganizationHandler(orgService)

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(userRepo, authRepo, db, rdb, hasher, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	api := router.Group("/api/v1")
	{
		api.GET("/health", handler.HealthCheck)

		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.GET("/verify-email", authHandler.VerifyEmail)
		}

		// Organization Routes
		orgGroup := api.Group("/organizations")

		// Apply authentication middleware to all org routes
		orgGroup.Use(middleware.Auth(jwtManager, rdb))
		{
			// Bulk operations (Must come before /:id routes to avoid conflict)
			orgGroup.DELETE("/bulk", middleware.Permission("auth:organization:delete", db, rdb), orgHandler.BulkDelete)
			orgGroup.PUT("/bulk/restore", middleware.Permission("auth:organization:restore", db, rdb), orgHandler.BulkRestore)

			// Deleted entities fetching
			orgGroup.GET("/deleted", middleware.Permission("auth:organization:read", db, rdb), orgHandler.GetAllDeleted)
			orgGroup.GET("/deleted/:id", middleware.Permission("auth:organization:read", db, rdb), orgHandler.GetDeletedByID)

			// Example: Require specific permissions for creating and deleting
			orgGroup.POST("", middleware.Permission("auth:organization:create", db, rdb), orgHandler.Create)
			orgGroup.DELETE("/:id", middleware.Permission("auth:organization:delete", db, rdb), orgHandler.Delete)

			// General authenticated endpoints
			orgGroup.GET("", middleware.Permission("auth:organization:read", db, rdb), orgHandler.GetAll)
			orgGroup.GET("/:id", middleware.Permission("auth:organization:read", db, rdb), orgHandler.GetByID)
			orgGroup.PUT("/:id", middleware.Permission("auth:organization:update", db, rdb), orgHandler.Update)
			orgGroup.PUT("/:id/restore", middleware.Permission("auth:organization:restore", db, rdb), orgHandler.Restore)
		}
	}
}
