package config

import (
	"log"
	"os"
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
		log.Fatal("GCP_PROJECT_ID is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
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
