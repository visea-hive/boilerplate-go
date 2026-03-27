package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/visea-hive/auth-core/pkg/messages"
)

const LangKey = "lang"

// Lang parses the Accept-Language header and stores the resolved language code
// in the gin context under the key "lang". Handlers retrieve it with Lang(c).
func Lang() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(LangKey, messages.ParseLang(c.GetHeader("Accept-Language")))
		c.Next()
	}
}

// GetLang retrieves the language code set by the Lang middleware.
// Falls back to messages.LangDefault if the middleware was not applied.
func GetLang(c *gin.Context) string {
	if lang, exists := c.Get(LangKey); exists {
		if s, ok := lang.(string); ok {
			return s
		}
	}
	return messages.LangDefault
}
