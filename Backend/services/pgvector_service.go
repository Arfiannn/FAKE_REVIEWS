package services

import (
	"BE_FAKE_REVIEW/models"
	"fmt"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type PGVectorService struct {
	DB       *gorm.DB
	Embedder EmbeddingService
}

func NewPGVectorService(db *gorm.DB, embedder EmbeddingService) PGVectorService {
	return PGVectorService{
		DB:       db,
		Embedder: embedder,
	}
}

func (s PGVectorService) ImportReviewsFromCSV(csvPath string) (int, error) {
	reviews, err := LoadReviewsFromCSV(csvPath)
	if err != nil {
		return 0, err
	}

	inserted := 0

	for _, review := range reviews {
		var existing models.Review

		err := s.DB.Where("clean_review = ?", review.CleanReview).First(&existing).Error
		if err == nil {
			continue
		}

		if err := s.DB.Create(&review).Error; err != nil {
			return inserted, err
		}

		inserted++
	}

	return inserted, nil
}

func (s PGVectorService) GenerateEmbeddingsForReviews(limit int) (int, error) {
	if limit <= 0 {
		limit = 50
	}

	var reviews []models.Review

	if err := s.DB.
		Where("embedding IS NULL").
		Limit(limit).
		Find(&reviews).Error; err != nil {
		return 0, err
	}

	updated := 0

	for _, review := range reviews {
		vector, err := s.Embedder.CreateEmbedding(review.CleanReview)
		if err != nil {
			fmt.Println("Gagal embedding review ID:", review.ID, err)
			continue
		}

		pgVector := pgvector.NewVector(vector)

		if err := s.DB.Exec(
			"UPDATE reviews SET embedding = ? WHERE id = ?",
			pgVector,
			review.ID,
		).Error; err != nil {
			return updated, err
		}

		updated++
		fmt.Println("Berhasil embedding review ID:", review.ID)
	}

	return updated, nil
}

func (s PGVectorService) SearchSimilarByText(query string, topK int) ([]models.SearchResult, error) {
	vector, err := s.Embedder.CreateEmbedding(query)
	if err != nil {
		return nil, err
	}

	return s.SearchSimilarByVector(vector, topK)
}

func (s PGVectorService) SearchSimilarByVector(vector []float32, topK int) ([]models.SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}

	pgVector := pgvector.NewVector(vector)

	var results []models.SearchResult

	err := s.DB.Raw(`
		SELECT
			id,
			product_name,
			rating,
			clean_review,
			label,
			label_code,
			1 - (embedding <=> ?) AS similarity
		FROM reviews
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> ?
		LIMIT ?
	`, pgVector, pgVector, topK).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
