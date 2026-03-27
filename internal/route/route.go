package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/visea-hive/auth-core/internal/handler"
)

// InitRoutes sets up all routes for the application.
func InitRoutes(router *gin.Engine, db *gorm.DB) {
	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/health", handler.HealthCheck)

		// Register route module here
		RegisterProductRoutes(apiV1, db)
	}
}
