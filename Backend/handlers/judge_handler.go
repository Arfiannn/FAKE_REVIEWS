package handlers

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/models"
	"BE_FAKE_REVIEW/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JudgeHandler struct {
	Config config.Config
	DB     *gorm.DB
}

func NewJudgeHandler(cfg config.Config, db *gorm.DB) JudgeHandler {
	return JudgeHandler{
		Config: cfg,
		DB:     db,
	}
}

func (h JudgeHandler) JudgeClassification(c *gin.Context) {
	var input models.JudgeRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
		})
		return
	}

	llmService := services.NewDeepSeekService(h.Config)
	judgeService := services.NewJudgeService(llmService)

	result, err := judgeService.JudgeClassification(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "LLM-as-a-Judge berhasil menilai hasil klasifikasi",
		"data":    result,
	})
}
