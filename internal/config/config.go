package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	PASETO      PasetoConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

type PasetoConfig struct {
	SymmetricKey    string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Validate required environment variables
	symmetricKey := getEnv("PASETO_SYMMETRIC_KEY", "")
	if symmetricKey == "" {
		return nil, fmt.Errorf("PASETO_SYMMETRIC_KEY environment variable is required. Generate one with: openssl rand -hex 32")
	}

	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
			IdleTimeout:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "ecommerce"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		PASETO: PasetoConfig{
			SymmetricKey:    getEnv("PASETO_SYMMETRIC_KEY", ""),
			AccessDuration:  time.Duration(getEnvAsInt("PASETO_ACCESS_DURATION", 15)) * time.Minute,
			RefreshDuration: time.Duration(getEnvAsInt("PASETO_REFRESH_DURATION", 7)) * 24 * time.Hour,
		},
	}

	// Validate required fields
	if cfg.PASETO.SymmetricKey == "" {
		return nil, fmt.Errorf("PASETO_SYMMETRIC_KEY environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
