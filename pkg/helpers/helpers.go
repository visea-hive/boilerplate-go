package helpers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9\s-]`)
	multiSpaceRegex      = regexp.MustCompile(`[\s-]+`)
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

// ParseIDParam parses a numeric URL parameter as uint.
func ParseIDParam(c *gin.Context, param string) (uint, error) {
	idParam := c.Param(param)
	id, err := strconv.ParseUint(idParam, 10, 32)
	return uint(id), err
}

// ParseUUIDParam returns a UUID URL parameter as a string.
func ParseUUIDParam(c *gin.Context, param string) string {
	return c.Param(param)
}

// GenerateSlug converts a string into a URL-friendly slug.
// It lowercases, removes special characters, and collapses whitespace/hyphens.
func GenerateSlug(name string) string {
	lower := strings.ToLower(name)
	clean := nonAlphanumericRegex.ReplaceAllString(lower, "")
	slug := multiSpaceRegex.ReplaceAllString(clean, "-")
	return strings.Trim(slug, "-")
}

// Ptr returns a pointer to the given value. Useful for optional struct fields.
func Ptr[T any](v T) *T {
	return &v
}
