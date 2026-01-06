package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kk120306/cvwo-2026/backend/helpers"
)

// Config struct for holding all configurations
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
	JWT      JWTConfig
	AWS      AWSConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string // development or production
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL            string
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	SSLMode        string
	ChannelBinding string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// AWSConfig holds AWS configuration
type AWSConfig struct {
	AccessKeyID            string
	SecretAccessKey        string
	CloudFrontDistribution string
}

// loads and returns application configuration
func Load() *Config {
	// Load environment variables
	helpers.LoadEnvVariables()

	// Build database configuration
	dbConfig := DatabaseConfig{
		URL:            getEnv("DATABASE_URL", ""),
		Host:           getEnv("DB_HOST", ""),
		Port:           getEnv("DB_PORT", "5432"),
		User:           getEnv("DB_USER", ""),
		Password:       getEnv("DB_PASSWORD", ""),
		Name:           getEnv("DB_NAME", ""),
		SSLMode:        getEnv("DB_SSLMODE", "require"),
		ChannelBinding: getEnv("DB_CHANNEL_BINDING", "require"),
	}

	// for individual parts, have to build it into a URL
	if dbConfig.URL == "" && dbConfig.Host != "" {
		dbConfig.URL = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s channel_binding=%s",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.Port,
			dbConfig.SSLMode,
			dbConfig.ChannelBinding,
		)
	}

	// return pointer cus we dont want a copy
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: dbConfig,
		CORS: CORSConfig{
			AllowedOrigins:   []string{getEnv("FRONTEND_URL", "http://localhost:3000")},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
		JWT: JWTConfig{
			Secret: getEnv("SECRET", ""),
		},
		AWS: AWSConfig{
			AccessKeyID:            getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:        getEnv("AWS_SECRET_ACCESS_KEY", ""),
			CloudFrontDistribution: getEnv("CLOUDFRONT_DISTRIBUTION_ID", ""),
		},
	}
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		log.Fatal("Database configuration is required in the env file")
	}
	if c.JWT.Secret == "" {
		log.Fatal("SECRET is required for JWT authentication in the env file")
	}
	return nil
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
