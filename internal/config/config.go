package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Database
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBMaxOpenConns int
	DBMaxIdleConns int

	// JWT
	JWTSecret            string
	JWTExpiration        time.Duration
	JWTRefreshExpiration time.Duration

	// Server
	Port        string
	Environment string
	APIVersion  string

	// Security
	BCryptCost int

	// Logging
	LogLevel  string
	LogFormat string
}

func Load() *Config {
	return &Config{
		// Database
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "auth_db"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),

		// JWT
		JWTSecret:            getEnv("JWT_SECRET", "dev-jwt-secret-key-change-in-production"),
		JWTExpiration:        getEnvAsDuration("JWT_EXPIRATION", "24h"),
		JWTRefreshExpiration: getEnvAsDuration("JWT_REFRESH_EXPIRATION", "168h"),

		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENV", "development"),
		APIVersion:  getEnv("API_VERSION", "v1"),

		// Security
		BCryptCost: getEnvAsInt("BCRYPT_COST", 12),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "text"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
