package config

import (
	"fmt"
	"log"
	"net/mail"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Env *EnvConfig
)

type EnvConfig struct {
	Server struct {
		Port        string
		Environment string
	}
	DB struct {
		URI      string
		Database string
	}
}

// ValidationRule defines a validation function that returns an error if validation fails
type ValidationRule struct {
	Field   string
	Rule    func(value string) bool
	Message string
}

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		if GetEnv("GO_ENV") != "production" {
			log.Printf("Warning: .env file not found")
		}
	}

	Env = &EnvConfig{
		Server: struct {
			Port        string
			Environment string
		}{
			Port:        GetEnvOrDefault("PORT", "3005"),
			Environment: GetEnvOrDefault("GO_ENV", "development"),
		},
		DB: struct {
			URI      string
			Database string
		}{
			URI:      GetEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
			Database: GetEnvOrDefault("MONGO_DB", "events-api"),
		},
	}

	return validateConfig()
}

// validateConfig checks all required configuration values
func validateConfig() error {
	rules := []ValidationRule{
		// Server validation
		{
			Field:   "Server.Port",
			Rule:    func(v string) bool { return v != "" },
			Message: "server port is required",
		},

		// Database validation
		{
			Field:   "DB.URI",
			Rule:    func(v string) bool { return v != "" },
			Message: "database uri is required",
		},
		{
			Field:   "DB.Database",
			Rule:    func(v string) bool { return v != "" },
			Message: "database name is required",
		},
	}

	var errors []string
	for _, rule := range rules {
		value := getConfigValue(rule.Field)
		if !rule.Rule(value) {
			errors = append(errors, rule.Message)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// getConfigValue retrieves a configuration value using reflection based on the field path
func getConfigValue(fieldPath string) string {
	parts := strings.Split(fieldPath, ".")
	value := reflect.ValueOf(Env).Elem()

	for _, part := range parts {
		value = value.FieldByName(part)
	}

	return value.String()
}

// AddValidationRule allows adding custom validation rules
func AddValidationRule(field string, rule func(string) bool, message string) {
	customRules = append(customRules, ValidationRule{
		Field:   field,
		Rule:    rule,
		Message: message,
	})
}

// Custom validation rules that can be added by the application
var customRules []ValidationRule

// Custom validation helper functions
func IsValidPort(port string) bool {
	if port == "" {
		return false
	}
	portNum, err := strconv.Atoi(port)
	return err == nil && portNum > 0 && portNum <= 65535
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := GetEnv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
