package routes

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, cfg config.Config, db *gorm.DB) {
	reviewHandler := handlers.NewReviewHandler(cfg, db)
	ragHandler := handlers.NewRAGHandler(cfg, db)
	evaluationHandler := handlers.NewEvaluationHandler(cfg, db)
	judgeHandler := handlers.NewJudgeHandler(cfg, db)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Backend fake review berjalan",
			})
		})

		dataset := api.Group("/dataset")
		{
			dataset.GET("/info", reviewHandler.GetDatasetInfo)
			dataset.POST("/import", reviewHandler.ImportDataset)
		}

		embedding := api.Group("/embedding")
		{
			embedding.POST("/generate", reviewHandler.GenerateEmbedding)
		}

		rag := api.Group("/rag")
		{
			rag.POST("/search", ragHandler.SearchSimilar)
			rag.POST("/classify", ragHandler.ClassifyReview)
		}

		evaluation := api.Group("/evaluation")
		{
			evaluation.POST("/classification", evaluationHandler.EvaluateClassification)
			evaluation.POST("/retrieval", evaluationHandler.EvaluateRetrieval)
		}

		judge := api.Group("/judge")
		{
			judge.POST("/classification", judgeHandler.JudgeClassification)
		}
	}
}
