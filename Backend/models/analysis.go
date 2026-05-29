package models

import "time"

type ReviewAnalysis struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	InputReview     string    `gorm:"type:text;not null" json:"input_review"`
	HyDEDocument    string    `gorm:"type:text" json:"hyde_document"`
	PredictionLabel string    `gorm:"type:text" json:"prediction_label"`
	Confidence      string    `gorm:"type:text" json:"confidence"`
	ConfidenceScore int       `json:"confidence_score"`
	Reasoning       string    `gorm:"type:text" json:"reasoning"`
	RawAnswer       string    `gorm:"type:text" json:"raw_answer"`
	CreatedAt       time.Time `json:"created_at"`
}
