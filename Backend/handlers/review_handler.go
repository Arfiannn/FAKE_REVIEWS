package handlers

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/models"
	"BE_FAKE_REVIEW/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	Config config.Config
	DB     *gorm.DB
}

func NewReviewHandler(cfg config.Config, db *gorm.DB) ReviewHandler {
	return ReviewHandler{
		Config: cfg,
		DB:     db,
	}
}

func (h ReviewHandler) GetDatasetInfo(c *gin.Context) {
	var total int64
	var embedded int64
	var asli int64
	var palsu int64

	h.DB.Model(&models.Review{}).Count(&total)
	h.DB.Model(&models.Review{}).Where("embedding IS NOT NULL").Count(&embedded)
	h.DB.Model(&models.Review{}).Where("label_code = ?", 1).Count(&asli)
	h.DB.Model(&models.Review{}).Where("label_code = ?", 0).Count(&palsu)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Informasi dataset berhasil diambil",
		"data": gin.H{
			"total_data":        total,
			"total_embedding":   embedded,
			"total_belum_embed": total - embedded,
			"label_asli":        asli,
			"label_palsu":       palsu,
		},
	})
}

func (h ReviewHandler) ImportDataset(c *gin.Context) {
	embedder := services.NewEmbeddingService(h.Config)
	vectorService := services.NewPGVectorService(h.DB, embedder)

	total, err := vectorService.ImportReviewsFromCSV(h.Config.CSVPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dataset berhasil diimport ke PostgreSQL",
		"data": gin.H{
			"inserted": total,
		},
	})
}

func (h ReviewHandler) GenerateEmbedding(c *gin.Context) {
	type Request struct {
		Limit int `json:"limit"`
	}

	var input Request
	_ = c.ShouldBindJSON(&input)

	if input.Limit <= 0 {
		input.Limit = 50
	}

	embedder := services.NewEmbeddingService(h.Config)
	vectorService := services.NewPGVectorService(h.DB, embedder)

	total, err := vectorService.GenerateEmbeddingsForReviews(input.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Generate embedding selesai",
		"data": gin.H{
			"updated": total,
		},
	})
}
