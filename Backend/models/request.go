package models

type SearchRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

type ClassifyRequest struct {
	Review string `json:"review" binding:"required"`
	TopK   int    `json:"top_k"`
}

type SearchResult struct {
	ID          uint    `json:"id"`
	ProductName string  `json:"product_name"`
	Rating      int     `json:"rating"`
	CleanReview string  `json:"clean_review"`
	Label       string  `json:"label"`
	LabelCode   int     `json:"label_code"`
	Similarity  float64 `json:"similarity"`
}
