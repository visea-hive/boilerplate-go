package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/visea-hive/auth-core/pkg/logger"
	"github.com/visea-hive/auth-core/pkg/notifier"
)

// InitLogger configures the global slog logger based on environment.
func InitLogger(env string) {
	var handler slog.Handler

	if env == "production" {
		// JSON format for production (easy to parse by log aggregators)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Human-readable text format for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	slog.SetDefault(slog.New(handler))
}

// InitTimezone sets the local timezone for the application.
func InitTimezone(timezone string) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Error("Failed to load timezone, using UTC", "timezone", timezone, "error", err)
		return
	}

	time.Local = loc
	slog.Info("Timezone configured", "timezone", timezone)
}

// InitCORS configures CORS middleware on the given Gin engine.
func InitCORS(router *gin.Engine, cfg *CORSConfig) {
	allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// InitNotifier creates the appropriate notifier from config, wires it into
// the unified logger (slog + notification in one call), and returns the logger.
func InitNotifier(cfg *NotificationConfig) *logger.Logger {
	var n notifier.Notifier

	if cfg.Enabled && cfg.WebhookURL != "" {
		n = notifier.NewWebhookNotifier(cfg.WebhookURL, cfg.Provider)
		slog.Info("Notification enabled", "provider", cfg.Provider)
	} else {
		n = notifier.NewNoOpNotifier()
		slog.Info("Notification disabled (using noop)")
	}

	log := logger.New(notifier.NewAsync(n))
	logger.SetDefault(log)
	return log
}

// InitRedis creates and validates a Redis client connection.
func InitRedis(cfg *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return rdb, nil
}
