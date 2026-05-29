package services

import (
	"BE_FAKE_REVIEW/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type EvaluationService struct {
	DB         *gorm.DB
	RAGService RAGService
}

func NewEvaluationService(db *gorm.DB, ragService RAGService) EvaluationService {
	return EvaluationService{
		DB:         db,
		RAGService: ragService,
	}
}

func (s EvaluationService) EvaluateClassification(limit int, topK int) (models.EvaluationResult, error) {
	if limit <= 0 {
		limit = 20
	}

	if topK <= 0 {
		topK = 5
	}

	// Ambil data seimbang: 50% Asli dan 50% Palsu
	halfLimit := limit / 2
	if halfLimit <= 0 {
		halfLimit = 1
	}

	var asliReviews []models.Review
	var palsuReviews []models.Review

	err := s.DB.
		Where("clean_review IS NOT NULL").
		Where("clean_review <> ''").
		Where("label_code = ?", 1).
		Order("RANDOM()").
		Limit(halfLimit).
		Find(&asliReviews).Error

	if err != nil {
		return models.EvaluationResult{}, err
	}

	err = s.DB.
		Where("clean_review IS NOT NULL").
		Where("clean_review <> ''").
		Where("label_code = ?", 0).
		Order("RANDOM()").
		Limit(limit - halfLimit).
		Find(&palsuReviews).Error

	if err != nil {
		return models.EvaluationResult{}, err
	}

	reviews := append(asliReviews, palsuReviews...)

	var items []models.EvaluationItem

	tp := 0
	tn := 0
	fp := 0
	fn := 0

	for i, review := range reviews {
		fmt.Printf("Evaluasi data ke-%d dari %d | Review ID: %d | Label asli: %s\n",
			i+1,
			len(reviews),
			review.ID,
			review.Label,
		)

		analysis, _, err := s.RAGService.ClassifyReview(review.CleanReview, topK)
		if err != nil {
			fmt.Println("Gagal evaluasi review ID:", review.ID, err)
			continue
		}

		actualLabel := normalizeLabel(review.Label)
		predictedLabel := normalizeLabel(analysis.PredictionLabel)

		isCorrect := actualLabel == predictedLabel

		if actualLabel == "palsu" && predictedLabel == "palsu" {
			tp++
		} else if actualLabel == "asli" && predictedLabel == "asli" {
			tn++
		} else if actualLabel == "asli" && predictedLabel == "palsu" {
			fp++
		} else if actualLabel == "palsu" && predictedLabel == "asli" {
			fn++
		}

		fmt.Println("Prediksi:", analysis.PredictionLabel, "| Benar:", isCorrect)

		items = append(items, models.EvaluationItem{
			ID:             review.ID,
			Review:         review.CleanReview,
			ActualLabel:    review.Label,
			PredictedLabel: analysis.PredictionLabel,
			IsCorrect:      isCorrect,
			Reasoning:      analysis.Reasoning,
		})
	}

	total := tp + tn + fp + fn

	accuracy := safeDivide(float64(tp+tn), float64(total))
	precision := safeDivide(float64(tp), float64(tp+fp))
	recall := safeDivide(float64(tp), float64(tp+fn))
	f1Score := safeDivide(2*precision*recall, precision+recall)

	result := models.EvaluationResult{
		TotalData: total,
		Accuracy:  accuracy,
		Precision: precision,
		Recall:    recall,
		F1Score:   f1Score,
		ConfusionMatrix: models.ConfusionMatrix{
			TP: tp,
			TN: tn,
			FP: fp,
			FN: fn,
		},
		Items: items,
	}

	return result, nil
}

func (s EvaluationService) EvaluateRetrieval(limit int, topK int) (models.RetrievalEvaluationResult, error) {
	if limit <= 0 {
		limit = 20
	}

	if topK <= 0 {
		topK = 5
	}

	halfLimit := limit / 2
	if halfLimit <= 0 {
		halfLimit = 1
	}

	var asliReviews []models.Review
	var palsuReviews []models.Review

	err := s.DB.
		Where("clean_review IS NOT NULL").
		Where("clean_review <> ''").
		Where("embedding IS NOT NULL").
		Where("label_code = ?", 1).
		Order("RANDOM()").
		Limit(halfLimit).
		Find(&asliReviews).Error

	if err != nil {
		return models.RetrievalEvaluationResult{}, err
	}

	err = s.DB.
		Where("clean_review IS NOT NULL").
		Where("clean_review <> ''").
		Where("embedding IS NOT NULL").
		Where("label_code = ?", 0).
		Order("RANDOM()").
		Limit(limit - halfLimit).
		Find(&palsuReviews).Error

	if err != nil {
		return models.RetrievalEvaluationResult{}, err
	}

	reviews := append(asliReviews, palsuReviews...)

	var items []models.RetrievalEvaluationItem

	totalPrecision := 0.0
	totalRecall := 0.0
	totalMRR := 0.0

	for i, review := range reviews {
		fmt.Printf("Evaluasi retrieval ke-%d dari %d | Review ID: %d | Label: %s\n",
			i+1,
			len(reviews),
			review.ID,
			review.Label,
		)

		results, err := s.RAGService.VectorDB.SearchSimilarByText(review.CleanReview, topK)
		if err != nil {
			fmt.Println("Gagal retrieval review ID:", review.ID, err)
			continue
		}

		queryLabel := normalizeLabel(review.Label)

		relevantCount := 0
		firstRelevantRank := 0

		for rank, result := range results {
			resultLabel := normalizeLabel(result.Label)

			// Hindari menghitung dokumen yang sama persis dengan query sebagai relevan
			if result.ID == review.ID {
				continue
			}

			if resultLabel == queryLabel {
				relevantCount++

				if firstRelevantRank == 0 {
					firstRelevantRank = rank + 1
				}
			}
		}

		precisionAtK := safeDivide(float64(relevantCount), float64(topK))
		recallAtK := safeDivide(float64(relevantCount), float64(topK))

		reciprocalRank := 0.0
		if firstRelevantRank > 0 {
			reciprocalRank = 1.0 / float64(firstRelevantRank)
		}

		totalPrecision += precisionAtK
		totalRecall += recallAtK
		totalMRR += reciprocalRank

		items = append(items, models.RetrievalEvaluationItem{
			ID:              review.ID,
			QueryReview:     review.CleanReview,
			QueryLabel:      review.Label,
			RelevantCount:   relevantCount,
			PrecisionAtK:    precisionAtK,
			RecallAtK:       recallAtK,
			ReciprocalRank:  reciprocalRank,
			RetrievedReviews: results,
		})
	}

	totalData := len(items)

	result := models.RetrievalEvaluationResult{
		TotalData:    totalData,
		TopK:         topK,
		PrecisionAtK: safeDivide(totalPrecision, float64(totalData)),
		RecallAtK:    safeDivide(totalRecall, float64(totalData)),
		MRR:          safeDivide(totalMRR, float64(totalData)),
		Items:        items,
	}

	return result, nil
}

func normalizeLabel(label string) string {
	label = strings.ToLower(strings.TrimSpace(label))

	if strings.Contains(label, "palsu") {
		return "palsu"
	}

	if strings.Contains(label, "asli") {
		return "asli"
	}

	if label == "0" {
		return "palsu"
	}

	if label == "1" {
		return "asli"
	}

	return label
}

func safeDivide(a float64, b float64) float64 {
	if b == 0 {
		return 0
	}

	return a / b
}