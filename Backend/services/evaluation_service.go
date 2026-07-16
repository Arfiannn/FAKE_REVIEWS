package services

import (
	"BE_FAKE_REVIEW/models"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EvaluationService struct {
	DB         *gorm.DB
	RAGService RAGService

	TestCSVPath              string
	RetrievalQueryPath       string
	RetrievalGroundTruthPath string
}

func NewEvaluationService(
	db *gorm.DB,
	ragService RAGService,
	testCSVPath string,
	retrievalQueryPath string,
	retrievalGroundTruthPath string,
) EvaluationService {
	return EvaluationService{
		DB:                       db,
		RAGService:               ragService,
		TestCSVPath:              testCSVPath,
		RetrievalQueryPath:       retrievalQueryPath,
		RetrievalGroundTruthPath: retrievalGroundTruthPath,
	}
}

func (s EvaluationService) EvaluateClassification(
	limit int,
	topK int,
) (models.EvaluationResult, error) {
	startTime := time.Now()

	if limit <= 0 {
		limit = 20
	}

	if topK <= 0 {
		topK = 5
	}

	allTestReviews, err := LoadEvaluationTestData(
		s.TestCSVPath,
	)
	if err != nil {
		return models.EvaluationResult{}, err
	}

	testReviews := selectBalancedTestReviews(
		allTestReviews,
		limit,
	)

	if len(testReviews) == 0 {
		return models.EvaluationResult{},
			fmt.Errorf("data test tidak tersedia")
	}

	items := make(
		[]models.EvaluationItem,
		0,
		len(testReviews),
	)

	tp := 0
	tn := 0
	fp := 0
	fn := 0

	for index, review := range testReviews {
		fmt.Printf(
			"Evaluasi klasifikasi ke-%d dari %d | Label aktual: %s\n",
			index+1,
			len(testReviews),
			review.Label,
		)

		analysis, _, err := s.RAGService.ClassifyReview(
			review.Review,
			topK,
		)

		if err != nil {
			fmt.Printf(
				"Gagal evaluasi data test ID %d: %v\n",
				review.ID,
				err,
			)
			continue
		}

		actualLabel := normalizeLabel(
			review.Label,
		)

		predictedLabel := normalizeLabel(
			analysis.PredictionLabel,
		)

		if actualLabel != "asli" &&
			actualLabel != "palsu" {
			fmt.Printf(
				"Label aktual tidak dikenal: %s\n",
				review.Label,
			)
			continue
		}

		if predictedLabel != "asli" &&
			predictedLabel != "palsu" {
			fmt.Printf(
				"Prediksi tidak dikenal: %s\n",
				analysis.PredictionLabel,
			)
			continue
		}

		isCorrect := actualLabel == predictedLabel

		// Palsu dianggap kelas positif
		if actualLabel == "palsu" &&
			predictedLabel == "palsu" {
			tp++
		} else if actualLabel == "asli" &&
			predictedLabel == "asli" {
			tn++
		} else if actualLabel == "asli" &&
			predictedLabel == "palsu" {
			fp++
		} else if actualLabel == "palsu" &&
			predictedLabel == "asli" {
			fn++
		}

		fmt.Printf(
			"Prediksi: %s | Aktual: %s | Benar: %t\n",
			analysis.PredictionLabel,
			review.Label,
			isCorrect,
		)

		items = append(
			items,
			models.EvaluationItem{
				ID:             review.ID,
				Review:         review.Review,
				ActualLabel:    review.Label,
				PredictedLabel: analysis.PredictionLabel,
				IsCorrect:      isCorrect,
				Reasoning:      analysis.Reasoning,
			},
		)
	}

	total := tp + tn + fp + fn

	if total == 0 {
		return models.EvaluationResult{},
			fmt.Errorf(
				"tidak ada evaluasi klasifikasi yang berhasil",
			)
	}

	accuracy := safeDivide(
		float64(tp+tn),
		float64(total),
	)

	precision := safeDivide(
		float64(tp),
		float64(tp+fp),
	)

	recall := safeDivide(
		float64(tp),
		float64(tp+fn),
	)

	f1Score := safeDivide(
		2*precision*recall,
		precision+recall,
	)

	executionTime := time.Since(startTime).Seconds()

	return models.EvaluationResult{
		TotalData:            total,
		Accuracy:             accuracy,
		Precision:            precision,
		Recall:               recall,
		F1Score:              f1Score,
		ExecutionTimeSeconds: executionTime,
		ConfusionMatrix: models.ConfusionMatrix{
			TP: tp,
			TN: tn,
			FP: fp,
			FN: fn,
		},
		Items: items,
	}, nil
}

func (s EvaluationService) EvaluateRetrieval(
	limit int,
	topK int,
) (models.RetrievalEvaluationResult, error) {
	startTime := time.Now()

	if limit <= 0 {
		limit = 100
	}

	if topK <= 0 {
		topK = 5
	}

	allQueries, err := LoadRetrievalQueries(
		s.RetrievalQueryPath,
	)
	if err != nil {
		return models.RetrievalEvaluationResult{}, err
	}

	groundTruth, err := LoadRetrievalGroundTruth(
		s.RetrievalGroundTruthPath,
	)
	if err != nil {
		return models.RetrievalEvaluationResult{}, err
	}

	if limit < len(allQueries) {
		allQueries = allQueries[:limit]
	}

	items := make(
		[]models.RetrievalEvaluationItem,
		0,
		len(allQueries),
	)

	totalPrecision := 0.0
	totalRecall := 0.0
	totalMRR := 0.0

	for index, query := range allQueries {
		relevantDocumentIDs, exists :=
			groundTruth[query.ID]

		if !exists ||
			len(relevantDocumentIDs) == 0 {
			fmt.Printf(
				"Ground truth query_id=%d tidak ditemukan\n",
				query.ID,
			)
			continue
		}

		fmt.Printf(
			"Evaluasi retrieval %d/%d | query_id=%d\n",
			index+1,
			len(allQueries),
			query.ID,
		)

		results, err := s.RAGService.VectorDB.
			SearchSimilarByText(
				query.Review,
				topK,
			)

		if err != nil {
			return models.RetrievalEvaluationResult{},
				fmt.Errorf(
					"gagal retrieval query_id=%d: %w",
					query.ID,
					err,
				)
		}

		relevantSet := createRelevantDocumentSet(
			relevantDocumentIDs,
		)

		if len(results) == 0 {
			items = append(
				items,
				models.RetrievalEvaluationItem{
					ID:               query.ID,
					QueryReview:      query.Review,
					QueryLabel:       query.Label,
					RelevantCount:    0,
					TotalRelevant:    len(relevantSet),
					PrecisionAtK:     0,
					RecallAtK:        0,
					ReciprocalRank:   0,
					RetrievedReviews: results,
				},
			)

			continue
		}

		relevantCount := 0
		firstRelevantRank := 0

		for rank, result := range results {
			// Menggunakan ID tabel reviews.
			if _, relevant :=
				relevantSet[result.ID]; !relevant {
				continue
			}

			relevantCount++

			if firstRelevantRank == 0 {
				firstRelevantRank = rank + 1
			}
		}

		precisionAtK := safeDivide(
			float64(relevantCount),
			float64(topK),
		)

		recallAtK := safeDivide(
			float64(relevantCount),
			float64(len(relevantSet)),
		)

		reciprocalRank := 0.0

		if firstRelevantRank > 0 {
			reciprocalRank =
				1.0 / float64(firstRelevantRank)
		}

		totalPrecision += precisionAtK
		totalRecall += recallAtK
		totalMRR += reciprocalRank

		items = append(
			items,
			models.RetrievalEvaluationItem{
				ID:               query.ID,
				QueryReview:      query.Review,
				QueryLabel:       query.Label,
				RelevantCount:    relevantCount,
				TotalRelevant:    len(relevantSet),
				PrecisionAtK:     precisionAtK,
				RecallAtK:        recallAtK,
				ReciprocalRank:   reciprocalRank,
				RetrievedReviews: results,
			},
		)
	}

	totalData := len(items)

	if totalData == 0 {
		return models.RetrievalEvaluationResult{},
			fmt.Errorf(
				"tidak ada evaluasi retrieval yang berhasil",
			)
	}

	averagePrecision := safeDivide(
		totalPrecision,
		float64(totalData),
	)

	averageRecall := safeDivide(
		totalRecall,
		float64(totalData),
	)

	averageMRR := safeDivide(
		totalMRR,
		float64(totalData),
	)

	return models.RetrievalEvaluationResult{
		TotalData: totalData,
		TopK:      topK,

		PrecisionAtK: math.Round(
			averagePrecision*1000000,
		) / 1000000,

		RecallAtK: math.Round(
			averageRecall*1000000,
		) / 1000000,

		MRR: math.Round(
			averageMRR*1000000,
		) / 1000000,

		ExecutionTimeSeconds: time.Since(startTime).Seconds(),

		Items: items,
	}, nil
}

func selectBalancedTestReviews(
	allReviews []models.EvaluationTestReview,
	limit int,
) []models.EvaluationTestReview {
	if limit <= 0 || limit > len(allReviews) {
		limit = len(allReviews)
	}

	targetAsli := limit / 2
	targetPalsu := limit - targetAsli

	asliReviews := make(
		[]models.EvaluationTestReview,
		0,
		targetAsli,
	)

	palsuReviews := make(
		[]models.EvaluationTestReview,
		0,
		targetPalsu,
	)

	for _, review := range allReviews {
		if review.LabelCode == 1 &&
			len(asliReviews) < targetAsli {
			asliReviews = append(
				asliReviews,
				review,
			)
		}

		if review.LabelCode == 0 &&
			len(palsuReviews) < targetPalsu {
			palsuReviews = append(
				palsuReviews,
				review,
			)
		}

		if len(asliReviews) >= targetAsli &&
			len(palsuReviews) >= targetPalsu {
			break
		}
	}

	result := make(
		[]models.EvaluationTestReview,
		0,
		limit,
	)

	selectedIDs := make(map[uint]bool)

	for _, review := range asliReviews {
		result = append(result, review)
		selectedIDs[review.ID] = true
	}

	for _, review := range palsuReviews {
		result = append(result, review)
		selectedIDs[review.ID] = true
	}

	if len(result) < limit {
		for _, review := range allReviews {
			if selectedIDs[review.ID] {
				continue
			}

			result = append(result, review)
			selectedIDs[review.ID] = true

			if len(result) >= limit {
				break
			}
		}
	}

	return result
}

func normalizeLabel(label string) string {
	label = strings.ToLower(
		strings.TrimSpace(label),
	)

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
