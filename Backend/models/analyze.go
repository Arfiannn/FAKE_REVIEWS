package models

type AnalyzeRequest struct {
	Input string `json:"input" binding:"required"`
	TopK  int    `json:"top_k"`
	Limit int    `json:"limit"`
}

type ShopeeReview struct {
	ProductName string `json:"product_name"`
	ShopName    string `json:"shop_name"`
	Username    string `json:"username"`
	Rating      int    `json:"rating"`
	Review      string `json:"review"`
	Date        string `json:"date"`
}

type ShopeeScraperResult struct {
	Success      bool           `json:"success"`
	TotalReviews int            `json:"total_reviews"`
	Reviews      []ShopeeReview `json:"reviews"`
	Error        string         `json:"error,omitempty"`
}

type AnalyzeReviewResult struct {
	ProductName      string         `json:"product_name,omitempty"`
	ShopName         string         `json:"shop_name,omitempty"`
	Username         string         `json:"username,omitempty"`
	Rating           int            `json:"rating,omitempty"`
	Date             string         `json:"date,omitempty"`
	RawReview        string         `json:"raw_review"`
	CleanReview      string         `json:"clean_review"`
	Analysis         ReviewAnalysis `json:"analysis"`
	Judge            JudgeResult    `json:"judge"`
	RetrievalResults []SearchResult `json:"retrieval_results"`
	Error            string         `json:"error,omitempty"`
}

type AnalyzeSummary struct {
	TotalReview     int     `json:"total_review"`
	TotalAsli       int     `json:"total_asli"`
	TotalPalsu      int     `json:"total_palsu"`
	PercentageAsli  float64 `json:"percentage_asli"`
	PercentagePalsu float64 `json:"percentage_palsu"`
	ValidJudge      int     `json:"valid_judge"`
	NeedReviewJudge int     `json:"need_review_judge"`
}

type AnalyzeResponse struct {
	Type       string                `json:"type"`
	Input      string                `json:"input"`
	ProductURL string                `json:"product_url,omitempty"`
	Summary    *AnalyzeSummary       `json:"summary,omitempty"`
	Result     *AnalyzeReviewResult  `json:"result,omitempty"`
	Results    []AnalyzeReviewResult `json:"results,omitempty"`
}
