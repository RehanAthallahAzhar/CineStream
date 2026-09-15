package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                  string
	AppEnv                   string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	DBSslMode                string
	DBMaxOpenConns           int
	DBMaxIdleConns           int
	DBConnMaxLifetimeMinutes time.Duration
	JWTSecret                string
	JWTExpirationHours       int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		AppPort:                  getEnv("APP_PORT", "8080"),
		AppEnv:                   getEnv("APP_ENV", "development"),
		DBHost:                   getEnv("DB_HOST", "localhost"),
		DBPort:                   getEnv("DB_PORT", "5432"),
		DBUser:                   getEnv("DB_USER", "postgres"),
		DBPassword:               getEnv("DB_PASSWORD", "postgres"),
		DBName:                   getEnv("DB_NAME", "cinestream_db"),
		DBSslMode:                getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns:           getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:           getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetimeMinutes: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 15)) * time.Minute,
		JWTSecret:                getEnv("JWT_SECRET", "super-secret-key-cinestream-2026"),
		JWTExpirationHours:       getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
