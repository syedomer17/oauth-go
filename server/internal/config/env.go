package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppPort       string
	MongoURI      string
	MongoDB       string
	RedisAddr     string
	RedisPassword string
	RedisDB       string

	JWTSecret string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string

	GithubClientID     string
	GithubClientSecret string
	GithubCallbackURL  string

	UpstashRedisRestURL   string
	UpstashRedisRestToken string

	EmailFrom string
	SMTPHost  string
	SMTPPort  string
	SMTPUser  string
	SMTPPass  string

	ClientURL string

	MaxFileSize string
	UploadDir   string
}

var Config Env

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Config = Env{
		AppPort: getEnv("APP_PORT", "8000"),

		MongoURI:      getEnv("MONGO_URI", ""),
		MongoDB:       getEnv("MONGO_DB", ""),
		RedisAddr:     getEnv("REDIS_ADDR", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),

		JWTSecret: getEnv("JWT_SECRET", ""),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleCallbackURL:  getEnv("GOOGLE_CALLBACK_URL", ""),

		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubCallbackURL:  getEnv("GITHUB_CALLBACK_URL", ""),

		UpstashRedisRestURL:   getEnv("UPSTASH_REDIS_REST_URL", ""),
		UpstashRedisRestToken: getEnv("UPSTASH_REDIS_REST_TOKEN", ""),

		EmailFrom: getEnv("EMAIL_FROM", ""),
		SMTPHost:  getEnv("SMTP_HOST", ""),
		SMTPPort:  getEnv("SMTP_PORT", ""),
		SMTPUser:  getEnv("SMTP_USER", ""),
		SMTPPass:  getEnv("SMTP_PASS", ""),

		ClientURL: getEnv("CLIENT_URL", "http://localhost:5173"),

		MaxFileSize: getEnv("MAX_FILE_SIZE", "5242880"),
		UploadDir:   getEnv("UPLOAD_DIR", "uploads"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
