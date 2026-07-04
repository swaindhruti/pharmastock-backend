package config

import (
	"os"

	goenv "github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	UploadDir string

	JWTSecret      string
	AdminUsername  string
	AdminPassword  string
	AdminEmail     string
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() (*Config, error) {
	err := goenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "pharmastock-db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		UploadDir:  getEnv("UPLOAD_DIR", "./uploads"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		AdminUsername:  getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "admin123"),
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@pharmastock.com"),
	}

	return config, nil
}
