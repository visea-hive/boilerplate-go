package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
type Config struct {
	App          AppConfig
	DB           DBConfig
	CORS         CORSConfig
	Notification NotificationConfig
	Hash         HashConfig
	JWT          JWTConfig
	Redis        RedisConfig
}

// RedisConfig holds Redis cache connection settings.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// NotificationConfig holds notification provider settings.
type NotificationConfig struct {
	Enabled    bool
	Provider   string // "slack", "mattermost", etc.
	WebhookURL string
}

// HashConfig holds password hashing settings.
type HashConfig struct {
	Secret string // HASH_SECRET env var — used as argon2id pepper
}

// JWTConfig holds JWT signing settings.
type JWTConfig struct {
	Secret        string // JWT_SECRET env var
	AccessTTLMins int    // JWT_ACCESS_TTL_MINUTES (default: 15)
}

// CORSConfig holds CORS settings.
type CORSConfig struct {
	AllowedOrigins string
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Env      string
	Port     string
	Timezone string
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SSLMode  string
	Timezone string
}

// Load reads environment variables and returns a Config instance.
func Load() *Config {
	return &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "local"),
			Port:     getEnv("APP_PORT", "8080"),
			Timezone: getEnv("APP_TIMEZONE", "Asia/Jakarta"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "sister"),
			Port:     getEnv("DB_PORT", "5432"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		},
		Notification: NotificationConfig{
			Enabled:    getEnv("NOTIFICATION_ENABLED", "false") == "true",
			Provider:   getEnv("NOTIFICATION_PROVIDER", "slack"),
			WebhookURL: getEnv("NOTIFICATION_WEBHOOK_URL", ""),
		},
		Hash: HashConfig{
			Secret: getEnv("HASH_SECRET", ""),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", ""),
			AccessTTLMins: getEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
	}
}

// DSN returns the PostgreSQL connection string.
func (db *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s search_path=public",
		db.Host, db.User, db.Password, db.Name, db.Port, db.SSLMode, db.Timezone,
	)
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt reads an integer environment variable or returns a default value.
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err == nil {
			return n
		}
	}
	return defaultValue
}
