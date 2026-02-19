package config

import "os"

// Config holds application configuration
type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	Port       string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "azh"),
		DBPassword: getEnv("DB_PASSWORD", "azh"),
		DBName:     getEnv("DB_NAME", "azh"),
		DBPort:     getEnv("DB_PORT", "5432"),
		Port:       getEnv("PORT", "8080"),
	}
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
