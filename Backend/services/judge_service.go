package services

import (
	"BE_FAKE_REVIEW/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type JudgeService struct {
	LLM DeepSeekService
}

func NewJudgeService(llm DeepSeekService) JudgeService {
	return JudgeService{
		LLM: llm,
	}
}

func (s JudgeService) JudgeClassification(input models.JudgeRequest) (models.JudgeResult, error) {
	context := buildJudgeRetrievalContext(input.RetrievalResults)

	prompt := fmt.Sprintf(`
Kamu adalah LLM-as-a-Judge untuk mengevaluasi hasil klasifikasi fake review e-commerce.

Tugas:
Nilai apakah hasil klasifikasi berikut sudah tepat berdasarkan ulasan, alasan prediksi, confidence, dan dokumen retrieval.

(Input Ulasan)
%s

(Hasil Prediksi Sistem)
Label Prediksi: %s
Confidence: %s
Confidence Score: %d
Alasan Prediksi: %s

(Dokumen Retrieval)
%s

(Kriteria Penilaian)
Pertimbangkan:
- Apakah label prediksi sesuai dengan isi ulasan
- Apakah alasan prediksi logis
- Apakah confidence score sesuai dengan kekuatan bukti
- Apakah dokumen retrieval mendukung prediksi
- Apakah ulasan memiliki detail pengalaman nyata atau hanya pujian umum

(Output)
Berikan output dengan format berikut:
Judge Score: angka 0 sampai 100
Judge Verdict: Valid/Tidak Valid
Judge Comment: komentar singkat terhadap hasil prediksi.
`,
		input.Review,
		input.PredictionLabel,
		input.Confidence,
		input.ConfidenceScore,
		input.Reasoning,
		context,
	)

	answer, err := s.LLM.Chat(prompt)
	if err != nil {
		return models.JudgeResult{}, err
	}

	result := models.JudgeResult{
		JudgeScore:   extractJudgeScore(answer),
		JudgeVerdict: extractJudgeValue(answer, "Judge Verdict:"),
		JudgeComment: extractJudgeValue(answer, "Judge Comment:"),
		RawAnswer:    answer,
	}

	return result, nil
}

func buildJudgeRetrievalContext(results []models.SearchResult) string {
	var builder strings.Builder

	for i, item := range results {
		builder.WriteString(fmt.Sprintf("Dokumen %d\n", i+1))
		builder.WriteString(fmt.Sprintf("Similarity: %.4f\n", item.Similarity))
		builder.WriteString(fmt.Sprintf("Review: %s\n", item.CleanReview))
		builder.WriteString(fmt.Sprintf("Label Pembanding: %s\n", item.Label))
		builder.WriteString("\n")
	}

	return builder.String()
}

func extractJudgeValue(text string, key string) string {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToLower(line), strings.ToLower(key)) {
			return strings.TrimSpace(strings.TrimPrefix(line, key))
		}
	}

	return ""
}

func extractJudgeScore(text string) int {
	value := extractJudgeValue(text, "Judge Score:")

	re := regexp.MustCompile(`\d+`)
	match := re.FindString(value)

	if match == "" {
		return 0
	}

	score, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}

	if score < 0 {
		return 0
	}

	if score > 100 {
		return 100
	}

	return score
}
