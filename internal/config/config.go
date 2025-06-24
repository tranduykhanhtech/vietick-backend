package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpiryHour int
	RefreshExpiryDay int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	accessExpiryHour, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRY_HOUR", "24"))
	refreshExpiryDay, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_DAY", "7"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "vietick"),
		},
		JWT: JWTConfig{
			AccessSecret:     getEnv("JWT_ACCESS_SECRET", "your-super-secret-jwt-access-key"),
			RefreshSecret:    getEnv("JWT_REFRESH_SECRET", "your-super-secret-jwt-refresh-key"),
			AccessExpiryHour: accessExpiryHour,
			RefreshExpiryDay: refreshExpiryDay,
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     smtpPort,
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@vietick.com"),
			FromName:     getEnv("FROM_NAME", "VietTick"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
