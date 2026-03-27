package helpers

import (
	"context"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/visea-hive/auth-core/pkg/jwt"
)

// UserInformationKey is the context key used by the auth middleware to store
// the parsed JWT claims. Use this key when setting/getting user information
// from both Gin and standard Go contexts.
const UserInformationKey = "userInformation"

// GetUserInformation extracts the authenticated user's information from a Gin
// context. Returns nil if the auth middleware has not run or the token was invalid.
func GetUserInformation(c *gin.Context) *jwtpkg.Claims {
	v, exists := c.Get(UserInformationKey)
	if !exists {
		return nil
	}
	userInfo, _ := v.(*jwtpkg.Claims)
	return userInfo
}

// GetUserInformationCtx extracts the authenticated user's information from a
// standard Go context (useful in service and repository layers).
// Returns nil if the context does not contain user information.
func GetUserInformationCtx(ctx context.Context) *jwtpkg.Claims {
	userInfo, _ := ctx.Value(UserInformationKey).(*jwtpkg.Claims)
	return userInfo
}
