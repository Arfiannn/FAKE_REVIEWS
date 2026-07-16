package handlers

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/models"
	"BE_FAKE_REVIEW/services"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyzeHandler struct {
	Config config.Config
	DB     *gorm.DB
}

func NewAnalyzeHandler(cfg config.Config, db *gorm.DB) AnalyzeHandler {
	return AnalyzeHandler{
		Config: cfg,
		DB:     db,
	}
}

func (h AnalyzeHandler) Analyze(c *gin.Context) {
	var input models.AnalyzeRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
		})
		return
	}

	input.Input = strings.TrimSpace(input.Input)

	if input.TopK <= 0 {
		input.TopK = 10
	}

	if input.Limit <= 0 {
		input.Limit = 1
	}

	embedder := services.NewEmbeddingService(h.Config)
	vectorService := services.NewPGVectorService(h.DB, embedder)
	llmService := services.NewDeepSeekService(h.Config)
	ragService := services.NewRAGService(h.DB, embedder, vectorService, llmService)
	judgeService := services.NewJudgeService(llmService)

	if services.IsShopeeLink(input.Input) {
		h.analyzeShopeeLink(c, input, ragService, judgeService)
		return
	}

	h.analyzeSingleReview(c, input, ragService, judgeService)
}

func (h AnalyzeHandler) analyzeSingleReview(
	c *gin.Context,
	input models.AnalyzeRequest,
	ragService services.RAGService,
	judgeService services.JudgeService,
) {
	cleanReview := services.PreprocessReviewText(input.Input)

	analysis, retrievalResults, err := ragService.ClassifyReview(cleanReview, input.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	judgeResult, err := judgeService.JudgeClassification(models.JudgeRequest{
		Review:           cleanReview,
		PredictionLabel:  analysis.PredictionLabel,
		Confidence:       analysis.Confidence,
		ConfidenceScore:  analysis.ConfidenceScore,
		Reasoning:        analysis.Reasoning,
		RetrievalResults: retrievalResults,
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Analisis review berhasil, tetapi LLM-as-a-Judge gagal dijalankan",
			"data": models.AnalyzeResponse{
				Type:  "single_review",
				Input: input.Input,
				Result: &models.AnalyzeReviewResult{
					RawReview:        input.Input,
					CleanReview:      cleanReview,
					Analysis:         analysis,
					RetrievalResults: retrievalResults,
					Error:            err.Error(),
				},
			},
		})
		return
	}

	result := models.AnalyzeReviewResult{
		RawReview:        input.Input,
		CleanReview:      cleanReview,
		Analysis:         analysis,
		Judge:            judgeResult,
		RetrievalResults: retrievalResults,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analisis review berhasil",
		"data": models.AnalyzeResponse{
			Type:   "single_review",
			Input:  input.Input,
			Result: &result,
		},
	})
}

func (h AnalyzeHandler) analyzeShopeeLink(
	c *gin.Context,
	input models.AnalyzeRequest,
	ragService services.RAGService,
	judgeService services.JudgeService,
) {
	scraperService := services.NewShopeeScraperService(h.Config)

	scrapeResult, err := scraperService.ScrapeReviews(input.Input, input.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var results []models.AnalyzeReviewResult

	totalAsli := 0
	totalPalsu := 0
	validJudge := 0
	needReviewJudge := 0

	for _, review := range scrapeResult.Reviews {
		cleanReview := services.PreprocessReviewText(review.Review)

		if cleanReview == "" {
			continue
		}

		analysis, retrievalResults, err := ragService.ClassifyReview(cleanReview, input.TopK)
		if err != nil {
			results = append(results, models.AnalyzeReviewResult{
				ProductName: review.ProductName,
				ShopName:    review.ShopName,
				Username:    review.Username,
				Rating:      review.Rating,
				Date:        review.Date,
				RawReview:   review.Review,
				CleanReview: cleanReview,
				Error:       err.Error(),
			})
			continue
		}

		judgeResult, judgeErr := judgeService.JudgeClassification(models.JudgeRequest{
			Review:           cleanReview,
			PredictionLabel:  analysis.PredictionLabel,
			Confidence:       analysis.Confidence,
			ConfidenceScore:  analysis.ConfidenceScore,
			Reasoning:        analysis.Reasoning,
			RetrievalResults: retrievalResults,
		})

		if strings.Contains(strings.ToLower(analysis.PredictionLabel), "asli") {
			totalAsli++
		}

		if strings.Contains(strings.ToLower(analysis.PredictionLabel), "palsu") {
			totalPalsu++
		}

		if strings.Contains(strings.ToLower(judgeResult.JudgeVerdict), "valid") &&
			!strings.Contains(strings.ToLower(judgeResult.JudgeVerdict), "tidak") {
			validJudge++
		}

		if strings.Contains(strings.ToLower(judgeResult.JudgeVerdict), "tidak") {
			needReviewJudge++
		}

		item := models.AnalyzeReviewResult{
			ProductName:      review.ProductName,
			ShopName:         review.ShopName,
			Username:         review.Username,
			Rating:           review.Rating,
			Date:             review.Date,
			RawReview:        review.Review,
			CleanReview:      cleanReview,
			Analysis:         analysis,
			Judge:            judgeResult,
			RetrievalResults: retrievalResults,
		}

		if judgeErr != nil {
			item.Error = judgeErr.Error()
		}

		results = append(results, item)
	}

	totalReview := len(results)

	summary := models.AnalyzeSummary{
		TotalReview:     totalReview,
		TotalAsli:       totalAsli,
		TotalPalsu:      totalPalsu,
		PercentageAsli:  safePercent(totalAsli, totalReview),
		PercentagePalsu: safePercent(totalPalsu, totalReview),
		ValidJudge:      validJudge,
		NeedReviewJudge: needReviewJudge,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analisis review dari link Shopee berhasil",
		"data": models.AnalyzeResponse{
			Type:       "shopee_link",
			Input:      input.Input,
			ProductURL: input.Input,
			Summary:    &summary,
			Results:    results,
		},
	})
}

func safePercent(value int, total int) float64 {
	if total == 0 {
		return 0
	}

	percentage := float64(value) / float64(total) * 100

	return math.Round(percentage)
}
