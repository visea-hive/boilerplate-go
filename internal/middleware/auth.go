package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/visea-hive/auth-core/pkg/helpers"
	jwtpkg "github.com/visea-hive/auth-core/pkg/jwt"
	"github.com/visea-hive/auth-core/pkg/messages"
)

// Auth returns a Gin middleware that validates the Bearer JWT token in the
// Authorization header and stores the parsed user information in the context.
//
// Inject the *jwtpkg.Manager created at startup (via config.InitJWT) so the
// secret is never read from the environment at request time.
//
// It also checks a Redis blacklist for fast session/user/org revocation without
// hitting the database.
func Auth(jwtManager *jwtpkg.Manager, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := messages.ParseLang(c.GetHeader("Accept-Language"))
		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		if raw == "" {
			slog.Warn("Auth: missing token", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": messages.Translate(lang, messages.ErrUnauthorized),
			})
			return
		}

		userInfo, err := jwtManager.Verify(raw)
		if err != nil {
			msg := messages.ErrTokenInvalid
			if err == messages.ErrTokenExpired {
				msg = messages.ErrTokenExpired
				slog.Info("Auth: expired token", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			} else {
				slog.Warn("Auth: invalid token", "ip", c.ClientIP(), "path", c.Request.URL.Path, "error", err)
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": messages.Translate(lang, msg),
			})
			return
		}

		// Fast Redis Blacklist Check (Zero DB queries)
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		keys := []string{
			"blacklist:session:" + userInfo.SessionID,
			"blacklist:user:" + userInfo.UserUUID,
			"blacklist:org_access:" + userInfo.UserUUID + ":" + fmt.Sprint(userInfo.ActiveOrgID),
		}

		results, err := rdb.MGet(ctx, keys...).Result()
		if err == nil {
			for _, val := range results {
				if val != nil {
					// Token or user scope was explicitly revoked
					slog.Warn("Auth: token blacklisted", "user_uuid", userInfo.UserUUID, "session", userInfo.SessionID, "ip", c.ClientIP())
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": messages.Translate(lang, messages.ErrUnauthorized),
					})
					return
				}
			}
		} else {
			slog.Error("Auth: failed to check redis blacklist", "error", err)
			// Decide if we should fail open or closed here. Opting to fail open so a
			// Redis hiccup doesn't bring down the whole app, given token signature is valid.
		}

		// Store the full typed claims
		c.Set(helpers.UserInformationKey, userInfo)
		c.Next()
	}
}
