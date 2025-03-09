package config

import (
	"log"
	"os"
	"strconv"
)

// Config stores all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

// ServerConfig stores server related configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig stores database related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// StorageConfig stores file storage related configuration
type StorageConfig struct {
	Type           string // "s3", "local", etc.
	Endpoint       string
	Region         string
	Bucket         string
	AccessKey      string
	SecretKey      string
	ForcePathStyle bool
}

// AuthConfig stores authentication related configuration
type AuthConfig struct {
	JWTSecret     string
	TokenExpiry   int // in hours
	RefreshExpiry int // in days
}

// Load returns a Config struct filled with values from the environment
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "easy_storage"),
		},
		Storage: StorageConfig{
			Type:           getEnv("STORAGE_TYPE", "s3"),
			Endpoint:       getEnv("STORAGE_ENDPOINT", "http://localhost:9000"),
			Region:         getEnv("STORAGE_REGION", "us-east-1"),
			Bucket:         getEnv("STORAGE_BUCKET", "easy-storage"),
			AccessKey:      getEnv("STORAGE_ACCESS_KEY", "minio"),
			SecretKey:      getEnv("STORAGE_SECRET_KEY", "minio123"),
			ForcePathStyle: getEnvAsBool("STORAGE_FORCE_PATH_STYLE", true),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
			TokenExpiry:   getEnvAsInt("TOKEN_EXPIRY", 24),
			RefreshExpiry: getEnvAsInt("REFRESH_EXPIRY", 7),
		},
	}
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	log.Printf("Fetching environment variable: %s", key)
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Fetching environment value: %s", value)

		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
