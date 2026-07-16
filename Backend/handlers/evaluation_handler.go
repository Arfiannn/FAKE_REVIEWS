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

func NewEvaluationHandler(
	cfg config.Config,
	db *gorm.DB,
) EvaluationHandler {
	return EvaluationHandler{
		Config: cfg,
		DB:     db,
	}
}

// =====================================================
// EVALUASI KLASIFIKASI
// =====================================================

func (h EvaluationHandler) EvaluateClassification(
	c *gin.Context,
) {
	var input models.EvaluationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": "Input tidak valid",
			},
		)
		return
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}

	if input.TopK <= 0 {
		input.TopK = 5
	}

	evaluationService := h.createEvaluationService()

	result, err := evaluationService.
		EvaluateClassification(
			input.Limit,
			input.TopK,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Evaluasi klasifikasi berhasil",
			"data":    result,
		},
	)
}

// =====================================================
// MEMBUAT KANDIDAT ANOTASI RETRIEVAL
// =====================================================

func (h EvaluationHandler) GenerateRetrievalAnnotation(
	c *gin.Context,
) {
	var input models.RetrievalAnnotationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": "Input tidak valid",
			},
		)
		return
	}

	if input.Limit <= 0 {
		input.Limit = 100
	}

	if input.CandidateK <= 0 {
		input.CandidateK = 20
	}

	evaluationService := h.createEvaluationService()

	processedQueries, totalRows, err :=
		evaluationService.GenerateRetrievalAnnotation(
			h.Config.RetrievalQueryPath,
			h.Config.RetrievalAnnotationPath,
			input.Limit,
			input.CandidateK,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "File anotasi retrieval berhasil dibuat",
			"data": gin.H{
				"requested_queries": input.Limit,
				"processed_queries": processedQueries,
				"candidate_k":       input.CandidateK,
				"total_rows":        totalRows,
				"output_path":       h.Config.RetrievalAnnotationPath,
			},
		},
	)
}

// =====================================================
// MEMBUAT GROUND TRUTH DARI HASIL ANOTASI
// =====================================================

func (h EvaluationHandler) GenerateRetrievalGroundTruth(
	c *gin.Context,
) {
	totalQueries, totalRelevantDocuments, err :=
		services.GenerateRetrievalGroundTruth(
			h.Config.RetrievalAnnotationPath,
			h.Config.RetrievalGroundTruthPath,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Ground truth retrieval berhasil dibuat",
			"data": gin.H{
				"total_queries":            totalQueries,
				"total_relevant_documents": totalRelevantDocuments,
				"output_path":              h.Config.RetrievalGroundTruthPath,
			},
		},
	)
}

// =====================================================
// EVALUASI RETRIEVAL DENGAN GROUND TRUTH
// =====================================================

func (h EvaluationHandler) EvaluateRetrieval(
	c *gin.Context,
) {
	var input models.RetrievalEvaluationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": "Input tidak valid",
			},
		)
		return
	}

	if input.Limit <= 0 {
		input.Limit = 100
	}

	if input.TopK <= 0 {
		input.TopK = 5
	}

	evaluationService := h.createEvaluationService()

	result, err := evaluationService.
		EvaluateRetrieval(
			input.Limit,
			input.TopK,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Evaluasi retrieval berhasil",
			"data":    result,
		},
	)
}

// =====================================================
// HELPER MEMBUAT EVALUATION SERVICE
// =====================================================

func (h EvaluationHandler) createEvaluationService() services.EvaluationService {
	embedder := services.NewEmbeddingService(
		h.Config,
	)

	vectorService := services.NewPGVectorService(
		h.DB,
		embedder,
	)

	llmService := services.NewDeepSeekService(
		h.Config,
	)

	ragService := services.NewRAGService(
		h.DB,
		embedder,
		vectorService,
		llmService,
	)

	return services.NewEvaluationService(
		h.DB,
		ragService,
		h.Config.TestCSVPath,
		h.Config.RetrievalQueryPath,
		h.Config.RetrievalGroundTruthPath,
	)
}
