package config

import (
	"os"

	"github.com/task-manager/task-service/pkg/logger"
)

// Config holds the application configuration
type Config struct {
	Port          string
	GCPProjectID  string
	JWTSecret     string
	JWTExpiration string
}

// Load reads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Port:          getEnv("PORT", "8081"),
		GCPProjectID:  getEnv("GCP_PROJECT_ID", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
	}

	if cfg.GCPProjectID == "" {
		logger.Error("GCP_PROJECT_ID is required")
		os.Exit(1)
	}

	if cfg.JWTSecret == "" {
		logger.Error("JWT_SECRET is required")
		os.Exit(1)
	}

	return cfg
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
