package config

import (
	"os"
	"log"

   "github.com/joho/godotenv"
)

type Env struct {
	AppPort  string
	MongoURI string
	MongoDB  string

	JWTSecret string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string

	GithubClientID     string
	GithubClientSecret string
	GithubCallbackURL  string
}

var Config Env
func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Config = Env{
		AppPort: getEnv("APP_PORT", "8000"),

		MongoURI: getEnv("MONGO_URI", ""),
		MongoDB:  getEnv("MONGO_DB", ""),

		JWTSecret: getEnv("JWT_SECRET", ""),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleCallbackURL:  getEnv("GOOGLE_CALLBACK_URL", ""),

		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubCallbackURL:  getEnv("GITHUB_CALLBACK_URL", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}