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

	CSVPath     string
	TestCSVPath string

	RetrievalQueryPath       string
	RetrievalAnnotationPath  string
	RetrievalGroundTruthPath string

	OllamaURL        string
	OllamaEmbedModel string

	DeepSeekAPIKey string
	DeepSeekModel  string

	ShopeeUserDataDir  string
	ShopeeHeadless     string
	ShopeeDefaultLimit string
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

		CSVPath: getEnv(
			"CSV_PATH",
			"data/reviews_train.csv",
		),

		TestCSVPath: getEnv(
			"TEST_CSV_PATH",
			"data/reviews_test.csv",
		),

		RetrievalQueryPath: getEnv(
			"RETrieval_QUERY_PATH",
			"data/retrieval_queries.csv",
		),

		RetrievalAnnotationPath: getEnv(
			"RETrieval_ANNOTATION_PATH",
			"data/retrieval_annotation.csv",
		),

		RetrievalGroundTruthPath: getEnv(
			"RETrieval_GROUND_TRUTH_PATH",
			"data/retrieval_ground_truth.csv",
		),

		OllamaURL: getEnv(
			"OLLAMA_URL",
			"http://localhost:11434",
		),

		OllamaEmbedModel: getEnv(
			"OLLAMA_EMBED_MODEL",
			"nomic-embed-text",
		),

		DeepSeekAPIKey: getEnv(
			"DEEPSEEK_API_KEY",
			"",
		),

		DeepSeekModel: getEnv(
			"DEEPSEEK_MODEL",
			"deepseek-chat",
		),

		ShopeeUserDataDir: getEnv(
			"SHOPEE_USER_DATA_DIR",
			"chrome-profile-go",
		),

		ShopeeHeadless: getEnv(
			"SHOPEE_HEADLESS",
			"false",
		),

		ShopeeDefaultLimit: getEnv(
			"SHOPEE_DEFAULT_LIMIT",
			"5",
		),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
