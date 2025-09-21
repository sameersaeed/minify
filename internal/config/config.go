package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	BaseURL     string
	JWTSecret   string
}

// Load reads environment variables (via .env) and returns a Config struct with defaults.
func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
	    DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres@localhost/minify?sslmode=disable"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
		JWTSecret:   getEnv("JWT_SECRET", "changeit"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
