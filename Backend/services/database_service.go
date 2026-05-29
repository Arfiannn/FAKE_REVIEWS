package services

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupDatabase(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS vector`).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Review{}, &models.ReviewAnalysis{}); err != nil {
		return err
	}

	db.Exec(`
		CREATE INDEX IF NOT EXISTS reviews_embedding_hnsw_idx
		ON reviews
		USING hnsw (embedding vector_cosine_ops)
	`)

	return nil
}
