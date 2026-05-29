package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	CSVPath string

	OllamaURL        string
	OllamaEmbedModel string

	DeepSeekAPIKey string
	DeepSeekModel  string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	return Config{
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "fake_review_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		CSVPath: getEnv("CSV_PATH", "data/reviews_labeled.csv"),

		OllamaURL:        getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaEmbedModel: getEnv("OLLAMA_EMBED_MODEL", "nomic-embed-text"),

		DeepSeekAPIKey: getEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekModel:  getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}