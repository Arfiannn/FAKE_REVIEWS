package handlers

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/models"
	"BE_FAKE_REVIEW/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RAGHandler struct {
	Config config.Config
	DB     *gorm.DB
}

func NewRAGHandler(cfg config.Config, db *gorm.DB) RAGHandler {
	return RAGHandler{
		Config: cfg,
		DB:     db,
	}
}

func (h RAGHandler) SearchSimilar(c *gin.Context) {
	var input models.SearchRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
		})
		return
	}

	if input.TopK <= 0 {
		input.TopK = 5
	}

	embedder := services.NewEmbeddingService(h.Config)
	vectorService := services.NewPGVectorService(h.DB, embedder)

	results, err := vectorService.SearchSimilarByText(input.Query, input.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Retrieval Top-K berhasil",
		"data": gin.H{
			"query":   input.Query,
			"top_k":   input.TopK,
			"results": results,
		},
	})
}

func (h RAGHandler) ClassifyReview(c *gin.Context) {
	var input models.ClassifyRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
		})
		return
	}

	if input.TopK <= 0 {
		input.TopK = 10
	}

	embedder := services.NewEmbeddingService(h.Config)
	vectorService := services.NewPGVectorService(h.DB, embedder)
	llmService := services.NewDeepSeekService(h.Config)
	ragService := services.NewRAGService(h.DB, embedder, vectorService, llmService)

	analysis, retrievalResults, err := ragService.ClassifyReview(input.Review, input.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	judgeService := services.NewJudgeService(llmService)

	judgeResult, err := judgeService.JudgeClassification(models.JudgeRequest{
		Review:           input.Review,
		PredictionLabel:  analysis.PredictionLabel,
		Confidence:       analysis.Confidence,
		ConfidenceScore:  analysis.ConfidenceScore,
		Reasoning:        analysis.Reasoning,
		RetrievalResults: retrievalResults,
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Klasifikasi review berhasil, tetapi LLM-as-a-Judge gagal dijalankan",
			"data": gin.H{
				"analysis":          analysis,
				"retrieval_results": retrievalResults,
				"judge":             nil,
				"judge_error":       err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Klasifikasi review dan validasi LLM-as-a-Judge berhasil",
		"data": gin.H{
			"analysis":          analysis,
			"retrieval_results": retrievalResults,
			"judge":             judgeResult,
		},
	})
}
