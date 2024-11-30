package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDB struct {
		URI      string
		Database string
	}
	Server struct {
		Port        string
		Environment string
	}
}

var AppConfig Config

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("GO_ENV") != "production" {
			log.Printf("Warning: .env file not found")
		}
	}

	AppConfig = Config{
		MongoDB: struct {
			URI      string
			Database string
		}{
			URI:      getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnvOrDefault("MONGO_DB", "events-api"),
		},
		Server: struct {
			Port        string
			Environment string
		}{
			Port:        getEnvOrDefault("PORT", "3005"),
			Environment: getEnvOrDefault("GO_ENV", "development"),
		},
	}

	return validateConfig()
}

func validateConfig() error {
	var errors []string

	// Validate MongoDB configuration
	if AppConfig.MongoDB.URI == "" {
		errors = append(errors, "MONGO_URI is required")
	}
	if AppConfig.MongoDB.Database == "" {
		errors = append(errors, "MONGO_DB is required")
	}

	// Validate Server configuration
	if AppConfig.Server.Port == "" {
		errors = append(errors, "PORT is required")
	} else {
		// Validate port number format
		if _, err := strconv.Atoi(AppConfig.Server.Port); err != nil {
			errors = append(errors, "PORT must be a valid number")
		}
	}

	if AppConfig.Server.Environment == "" {
		errors = append(errors, "GO_ENV is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
