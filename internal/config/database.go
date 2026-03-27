package config

import (
	"log/slog"
	"os"

	"github.com/visea-hive/auth-core/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance.
var DB *gorm.DB

// InitDB establishes a connection to the PostgreSQL database.
func InitDB(cfg *DBConfig) {
	var err error

	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: NewSlogGormLogger(logger.Info),
	})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("Database connected successfully")
}

// InitMigration runs auto-migration for all registered models.
func InitMigration() {
	err := DB.AutoMigrate(
		&model.Organization{},
		&model.User{},
		&model.Membership{},
		&model.OTPChallenge{},
		&model.PasswordResetToken{},
		&model.EmailVerification{},
		&model.Device{},
		&model.Session{},
		&model.RefreshToken{},
		&model.Service{},
		&model.Permission{},
		&model.Role{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.LoginAttempt{},
	)
	if err != nil {
		slog.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}

	slog.Info("Database migrated successfully")
}
