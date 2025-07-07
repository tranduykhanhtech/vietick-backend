package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
	CORS     CORSConfig
	Redis    RedisConfig
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

type CORSConfig struct {
	AllowedOrigins []string
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	accessExpiryHour, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRY_HOUR", "24"))
	refreshExpiryDay, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_DAY", "7"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

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
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS",
				[]string{
					"http://localhost:3000",
					"http://localhost:3001",
					"http://localhost:5173", // For vue/vite dev
					"https://vietick.com",
					"https://www.vietick.com",
					"https://app.vietick.com",
				},
				",",
			),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Username: getEnv("REDIS_USERNAME", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultVal []string, separator string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	return strings.Split(valStr, separator)
}

func (c *Config) GetRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     c.Redis.Addr,
		Username: c.Redis.Username,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	}
}
