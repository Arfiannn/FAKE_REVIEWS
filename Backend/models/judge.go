package models

type JudgeRequest struct {
	Review           string         `json:"review" binding:"required"`
	PredictionLabel  string         `json:"prediction_label" binding:"required"`
	Confidence       string         `json:"confidence"`
	ConfidenceScore  int            `json:"confidence_score"`
	Reasoning        string         `json:"reasoning"`
	RetrievalResults []SearchResult `json:"retrieval_results"`
}

type JudgeResult struct {
	JudgeScore   int    `json:"judge_score"`
	JudgeVerdict string `json:"judge_verdict"`
	JudgeComment string `json:"judge_comment"`
	RawAnswer    string `json:"raw_answer"`
}
