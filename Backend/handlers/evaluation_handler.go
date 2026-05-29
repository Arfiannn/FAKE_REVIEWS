package handlers

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/models"
	"BE_FAKE_REVIEW/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EvaluationHandler struct {
	Config config.Config
	DB     *gorm.DB
}

func NewEvaluationHandler(cfg config.Config, db *gorm.DB) EvaluationHandler {
	return EvaluationHandler{
		Config: cfg,
		DB:     db,
	}
}

func (h EvaluationHandler) EvaluateClassification(c *gin.Context) {
	var input models.EvaluationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
		})
		return
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}

	if input.TopK <= 0 {
		input.TopK = 5
	}

	embedder := services.NewEmbeddingService(h.Config)
	vectorService := services.NewPGVectorService(h.DB, embedder)
	llmService := services.NewDeepSeekService(h.Config)
	ragService := services.NewRAGService(h.DB, embedder, vectorService, llmService)

	evaluationService := services.NewEvaluationService(h.DB, ragService)

	result, err := evaluationService.EvaluateClassification(input.Limit, input.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Evaluasi klasifikasi berhasil",
		"data":    result,
	})
}

func (h EvaluationHandler) EvaluateRetrieval(c *gin.Context) {
	var input models.RetrievalEvaluationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
		})
		return
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}

	if input.TopK <= 0 {
		input.TopK = 5
	}

	embedder := services.NewEmbeddingService(h.Config)
	vectorService := services.NewPGVectorService(h.DB, embedder)
	llmService := services.NewDeepSeekService(h.Config)
	ragService := services.NewRAGService(h.DB, embedder, vectorService, llmService)

	evaluationService := services.NewEvaluationService(h.DB, ragService)

	result, err := evaluationService.EvaluateRetrieval(input.Limit, input.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Evaluasi retrieval berhasil",
		"data":    result,
	})
}
