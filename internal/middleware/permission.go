package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/visea-hive/auth-core/internal/repository"
	"github.com/visea-hive/auth-core/pkg/helpers"
	"github.com/visea-hive/auth-core/pkg/messages"
	"gorm.io/gorm"
)

// PermissionKey is the context key under which the matched permission full_key
// is stored after a successful permission check.
const PermissionKey = "permission"

// PermKey builds a permission full_key in the canonical format:
//
//	service:module:access
func PermKey(service, module, access string) string {
	return service + ":" + module + ":" + access
}

// Permission returns a middleware that checks whether the user's active role
// has the given permission full_key (format: "service:module:access").
//
// Permissions are checked against a Redis cache. On a cache miss, it queries
// the database and populates Redis.
//
// Must be placed after Auth middleware so that userInformation is already
// in the context.
func Permission(fullKey string, db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	permRepo := repository.NewPermissionRepository(db)

	return func(c *gin.Context) {
		lang := messages.ParseLang(c.GetHeader("Accept-Language"))

		// 1. Get user information from context (set by Auth middleware)
		userInfo := helpers.GetUserInformation(c)
		if userInfo == nil || userInfo.UserUUID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": messages.Translate(lang, messages.ErrUnauthorized),
			})
			return
		}

		// 2. Require an active role
		if userInfo.ActiveRoleID == 0 {
			slog.Warn("Permission: no active role",
				"user_uuid", userInfo.UserUUID,
				"path", c.Request.URL.Path,
				"required", fullKey,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": messages.Translate(lang, messages.ErrForbidden),
			})
			return
		}

		// 3. Check Redis for permissions
		ctx := c.Request.Context()
		cacheKey := fmt.Sprintf("role_perms:%d", userInfo.ActiveRoleID)

		perms, err := rdb.SMembers(ctx, cacheKey).Result()
		if err != nil || len(perms) == 0 {
			// Cache miss — load from DB
			perms = permRepo.LoadPermissionsFromDB(ctx, userInfo.ActiveRoleID)
			if len(perms) > 0 {
				// Convert to slice of interface{} for rdb.SAdd
				var permsIface []interface{}
				for _, p := range perms {
					permsIface = append(permsIface, p)
				}
				rdb.SAdd(ctx, cacheKey, permsIface...)
				rdb.Expire(ctx, cacheKey, 15*time.Minute)
			}
		}

		// 4. Verify permission
		if !hasPermission(perms, fullKey) {
			slog.Warn("Permission: access denied",
				"user_uuid", userInfo.UserUUID,
				"role_id", userInfo.ActiveRoleID,
				"role_name", userInfo.ActiveRoleName,
				"required", fullKey,
				"path", c.Request.URL.Path,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": messages.Translate(lang, messages.ErrForbidden),
			})
			return
		}

		// 5. Grant — store the matched key for optional downstream use
		c.Set(PermissionKey, fullKey)
		c.Next()
	}
}

// hasPermission reports whether key appears in the permissions slice.
func hasPermission(permissions []string, key string) bool {
	for _, p := range permissions {
		if p == key {
			return true
		}
	}
	return false
}
