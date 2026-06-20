package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBDriver      string
	DatabaseURL   string
	SQLitePath    string
	JWTSecret     string
	TokenTTLHours int
}

func Load() Config {
	if err := godotenv.Overload(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic("failed to load .env: " + err.Error())
	}

	return Config{
		Port:          getEnv("PORT", "8080"),
		DBDriver:      getEnv("DB_DRIVER", "sqlite"),
		DatabaseURL:   getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=hotel_management port=5432 sslmode=disable"),
		SQLitePath:    getEnv("SQLITE_PATH", "./hotel_management.db"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-hotel-key"),
		TokenTTLHours: 24,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
