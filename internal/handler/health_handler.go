package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visea-hive/auth-core/pkg/helpers"
)

// HealthCheck returns a simple health status response.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, helpers.SuccessResponse("Server is running", nil))
}
