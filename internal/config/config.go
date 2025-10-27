package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	SuperUser SuperUserConfig
	Session  SessionConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

// SuperUserConfig holds superuser configuration
type SuperUserConfig struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// SessionConfig holds session configuration
type SessionConfig struct {
	SecretKey      string
	CookieName     string
	CookieMaxAge   int // in seconds
	CookieSecure   bool
	CookieHttpOnly bool
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
		},
		SuperUser: SuperUserConfig{
			Email:     getEnv("SUPERUSER_EMAIL", "admin@angidi.com"),
			Password:  getEnv("SUPERUSER_PASSWORD", "Admin@123"),
			FirstName: getEnv("SUPERUSER_FIRST_NAME", "Super"),
			LastName:  getEnv("SUPERUSER_LAST_NAME", "Admin"),
		},
		Session: SessionConfig{
			SecretKey:      getEnv("SESSION_SECRET", "change-me-in-production"),
			CookieName:     getEnv("SESSION_COOKIE_NAME", "angidi_session"),
			CookieMaxAge:   getEnvAsInt("SESSION_COOKIE_MAX_AGE", 86400), // 24 hours
			CookieSecure:   getEnvAsBool("SESSION_COOKIE_SECURE", false),
			CookieHttpOnly: getEnvAsBool("SESSION_COOKIE_HTTP_ONLY", true),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
