package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	BaseURL     string
	FrontendURL string
	DatabaseURL string
	JWTSecret   string
}

// Load reads environment variables (via .env) and returns a Config struct with defaults.
func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres@localhost/minify?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET"),
	}
}

// Validate ensures the environment variables are set correctly for the program
func (c *Config) Validate() error {
	var errs []string

	if c.Port == "" {
		errs = append(errs, "PORT is required")
	} else if _, err := strconv.Atoi(c.Port); err != nil {
		errs = append(errs, "PORT must be a valid number")
	}

	if c.BaseURL == "" {
		errs = append(errs, "BASE_URL is required")
	}

	if c.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}

	if c.JWTSecret == "" || c.JWTSecret == "changeit" {
		errs = append(errs, "JWT_SECRET should be set to a secure random value (for example: openssl rand -base64 32)")
	}

	if len(errs) > 0 {
		return errors.New("config validation failed:\n  - " + strings.Join(errs, "\n  - "))
	}

	return nil
}

func getEnv(key string, defaultValue ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}
