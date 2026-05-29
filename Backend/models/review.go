package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type Review struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	ProductName string           `gorm:"type:text" json:"product_name"`
	ShopName    string           `gorm:"type:text" json:"shop_name"`
	Username    string           `gorm:"type:text" json:"username"`
	Rating      int              `json:"rating"`
	Review      string           `gorm:"type:text" json:"review"`
	CleanReview string           `gorm:"type:text;not null" json:"clean_review"`
	Label       string           `gorm:"type:text" json:"label"`
	LabelCode   int              `json:"label_code"`
	Embedding   *pgvector.Vector `gorm:"type:vector(768)" json:"-"`
	CreatedAt   time.Time        `json:"created_at"`
}
