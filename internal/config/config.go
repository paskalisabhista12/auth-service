package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

// Config holds all configuration values
type Config struct {
	AppPort       string
	DatabaseURL   string
	JwtSecret     string
	Environment   string
	RedisAddress  string
	RedisPassword string
}

// LoadConfig loads variables from .env into Config struct
func LoadConfig() Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found, using system environment variables")
	}

	config := Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		JwtSecret:     getEnv("JWT_SECRET", "defaultsecret"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		RedisAddress:  getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASS", ""),
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
