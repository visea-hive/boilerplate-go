package helpers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ContainsString checks if a string is present in a slice of strings.
func ContainsString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// ParseIDParam is a local helper mapping param
func ParseIDParam(c *gin.Context, param string) (uint, error) {
	idParam := c.Param(param)
	id, err := strconv.ParseUint(idParam, 10, 32)
	return uint(id), err
}

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
