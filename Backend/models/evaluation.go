package models

type EvaluationRequest struct {
	Limit int `json:"limit"`
	TopK  int `json:"top_k"`
}

type EvaluationTestReview struct {
	ID        uint   `json:"id"`
	Review    string `json:"review"`
	Label     string `json:"label"`
	LabelCode int    `json:"label_code"`
}

type EvaluationItem struct {
	ID             uint   `json:"id"`
	Review         string `json:"review"`
	ActualLabel    string `json:"actual_label"`
	PredictedLabel string `json:"predicted_label"`
	IsCorrect      bool   `json:"is_correct"`
	Reasoning      string `json:"reasoning"`
}

type ConfusionMatrix struct {
	TP int `json:"tp"`
	TN int `json:"tn"`
	FP int `json:"fp"`
	FN int `json:"fn"`
}

type EvaluationResult struct {
	TotalData int `json:"total_data"`

	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	F1Score   float64 `json:"f1_score"`

	ExecutionTimeSeconds float64 `json:"execution_time_seconds"`

	ConfusionMatrix ConfusionMatrix `json:"confusion_matrix"`

	Items []EvaluationItem `json:"items"`
}

type RetrievalEvaluationRequest struct {
	Limit int `json:"limit"`
	TopK  int `json:"top_k"`
}

type RetrievalAnnotationRequest struct {
	Limit      int `json:"limit"`
	CandidateK int `json:"candidate_k"`
}

type RetrievalEvaluationItem struct {
	ID            uint   `json:"id"`
	QueryReview   string `json:"query_review"`
	QueryLabel    string `json:"query_label"`
	RelevantCount int    `json:"relevant_count"`
	TotalRelevant int    `json:"total_relevant"`

	PrecisionAtK   float64 `json:"precision_at_k"`
	RecallAtK      float64 `json:"recall_at_k"`
	ReciprocalRank float64 `json:"reciprocal_rank"`

	RetrievedReviews []SearchResult `json:"retrieved_reviews"`
}

type RetrievalEvaluationResult struct {
	TotalData    int     `json:"total_data"`
	TopK         int     `json:"top_k"`
	PrecisionAtK float64 `json:"precision_at_k"`
	RecallAtK    float64 `json:"recall_at_k"`
	MRR          float64 `json:"mrr"`

	ExecutionTimeSeconds float64 `json:"execution_time_seconds"`

	Items []RetrievalEvaluationItem `json:"items"`
}
