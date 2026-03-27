package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/visea-hive/auth-core/internal/handler"
	"github.com/visea-hive/auth-core/internal/repository"
	"github.com/visea-hive/auth-core/internal/service"
)

func RegisterProductRoutes(router *gin.RouterGroup, db *gorm.DB) {
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	productRoutes := router.Group("/products")
	{
		productRoutes.GET("", productHandler.FindAll)
		productRoutes.GET("/datatable", productHandler.FindPaginated)
		productRoutes.GET("/:id", productHandler.FindByID)
		productRoutes.POST("", productHandler.Create)
		productRoutes.PUT("/:id", productHandler.Update)
		productRoutes.DELETE("/:id", productHandler.Delete)
	}
}
