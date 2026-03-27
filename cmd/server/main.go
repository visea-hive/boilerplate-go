package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/visea-hive/auth-core/internal/config"
	"github.com/visea-hive/auth-core/internal/route"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize structured logger
	config.InitLogger(cfg.App.Env)
	slog.Info("Starting application...", "env", cfg.App.Env)

	// Set local timezone
	config.InitTimezone(cfg.App.Timezone)

	// Set Gin mode based on environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	config.InitDB(&cfg.DB)

	// Initialize Redis
	rdb, err := config.InitRedis(&cfg.Redis)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	slog.Info("Redis connected successfully")

	// Run migrations (skip in local environment for efficiency)
	if cfg.App.Env != "local" {
		config.InitMigration()
	} else {
		slog.Info("Skipping auto-migration in local environment")
	}

	// Initialize notification service + unified logger
	log := config.InitNotifier(&cfg.Notification)

	// Setup Gin router with explicit middleware
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Initialize JWT Manager
	jwtManager := config.InitJWT(&cfg.JWT)

	// Initialize Hasher
	hasher := config.InitHasher(&cfg.Hash)

	// Initialize CORS
	config.InitCORS(router, &cfg.CORS)

	// Initialize Routes
	route.InitRoutes(router, rdb, config.DB, jwtManager, hasher)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Info("Server starting", "address", addr)
	if err := router.Run(addr); err != nil {
		log.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
