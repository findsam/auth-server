package config

import (
	"os"

	t "github.com/findsam/food-server/types"
	_ "github.com/joho/godotenv/autoload"
)

var Envs = initConfig()

func initConfig() t.Config {
	return t.Config{
		Env:              getEnv("ENV", "development"),
		Port:             getEnv("PORT", "8080"),
		MongoURI:         getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		PublicURL:        getEnv("PUBLIC_URL", "http://localhost:3000"),
		JWTSecret:        getEnv("JWT_SECRET", "JWT secret is required"),
		APIKey:           getEnv("API_KEY", "API Key is required"),
		ChatGPTSecretKey: getEnv("CHATGPT_SECRET_KEY", "ChatGPT API Key is required"),
		ChatGPTURL:       getEnv("CHATGPT_URL", "ChatGPT Url is required"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
