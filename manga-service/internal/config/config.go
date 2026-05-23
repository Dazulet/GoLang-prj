package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	AppEnv     string
	AppBaseURL string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string

	AuthServiceURL string

	UploadDir       string
	MaxUploadSizeMB int64
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading environment directly")
	}
	return &Config{
		AppPort:    getEnv("APP_PORT", "8082"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8082"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "manga_user"),
		DBPassword: getEnv("DB_PASSWORD", "manga_secret"),
		DBName:     getEnv("DB_NAME", "manga_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:      getEnv("JWT_SECRET", "fallback-secret"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),

		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB: int64(getEnvInt("MAX_UPLOAD_SIZE_MB", 10)),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
