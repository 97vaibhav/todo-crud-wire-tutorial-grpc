package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	GRPCPort   string
	JWTSecret  string // used to sign and verify JWT tokens
}

// Load reads .env file if present, then falls back to real environment variables.
// Wire calls this once at startup and the result is injected everywhere it's needed.
func Load() (*Config, error) {
	// godotenv.Load is intentionally ignored — if no .env file, that's fine.
	// The real environment (e.g. Docker env vars) takes precedence anyway.
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "tododb"),
		GRPCPort:   getEnv("GRPC_PORT", "50051"),
		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
