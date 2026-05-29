package services

import (
	"BE_FAKE_REVIEW/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type RAGService struct {
	DB       *gorm.DB
	Embedder EmbeddingService
	VectorDB PGVectorService
	LLM      DeepSeekService
}

func NewRAGService(db *gorm.DB, embedder EmbeddingService, vectorDB PGVectorService, llm DeepSeekService) RAGService {
	return RAGService{
		DB:       db,
		Embedder: embedder,
		VectorDB: vectorDB,
		LLM:      llm,
	}
}

func (s RAGService) GenerateHyDE(review string) (string, error) {
	prompt := fmt.Sprintf(`
Buat satu dokumen hipotetik singkat dalam bahasa Indonesia berdasarkan ulasan berikut.

Ulasan:
%s

Tujuan dokumen hipotetik ini adalah membantu sistem mencari ulasan lain yang mirip secara makna.
Jangan menentukan label Asli atau Palsu.
Tulis hanya dokumen hipotetiknya.
`, review)

	return s.LLM.Chat(prompt)
}

func (s RAGService) ClassifyReview(review string, topK int) (models.ReviewAnalysis, []models.SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}

	hydeDocument, err := s.GenerateHyDE(review)
	if err != nil {
		return models.ReviewAnalysis{}, nil, err
	}

	hydeEmbedding, err := s.Embedder.CreateEmbedding(hydeDocument)
	if err != nil {
		return models.ReviewAnalysis{}, nil, err
	}

	retrievalResults, err := s.VectorDB.SearchSimilarByVector(hydeEmbedding, topK)
	if err != nil {
		return models.ReviewAnalysis{}, nil, err
	}

	context := buildRetrievalContext(retrievalResults)
	prompt := buildClassificationPrompt(review, context)

	answer, err := s.LLM.Chat(prompt)
	if err != nil {
		return models.ReviewAnalysis{}, nil, err
	}

	label := extractValue(answer, "Label:")
	confidence := extractValue(answer, "Confidence:")
	confidenceScore := extractConfidenceScore(answer)
	reasoning := extractValue(answer, "Alasan:")

	analysis := models.ReviewAnalysis{
		InputReview:     review,
		HyDEDocument:    hydeDocument,
		PredictionLabel: label,
		Confidence:      confidence,
		ConfidenceScore: confidenceScore,
		Reasoning:       reasoning,
		RawAnswer:       answer,
	}

	if err := s.DB.Create(&analysis).Error; err != nil {
		return models.ReviewAnalysis{}, nil, err
	}

	return analysis, retrievalResults, nil
}

func buildRetrievalContext(results []models.SearchResult) string {
	var builder strings.Builder

	for i, item := range results {
		builder.WriteString(fmt.Sprintf("Dokumen %d\n", i+1))
		builder.WriteString(fmt.Sprintf("Similarity: %.4f\n", item.Similarity))
		builder.WriteString(fmt.Sprintf("Product: %s\n", item.ProductName))
		builder.WriteString(fmt.Sprintf("Rating: %d\n", item.Rating))
		builder.WriteString(fmt.Sprintf("Review: %s\n", item.CleanReview))
		builder.WriteString(fmt.Sprintf("Label pembanding: %s\n", item.Label))
		builder.WriteString("\n")
	}

	return builder.String()
}

func buildClassificationPrompt(review string, context string) string {
	return fmt.Sprintf(`
(Instruksi)
Analisis ulasan produk berikut dan tentukan apakah ulasan tersebut termasuk ulasan asli atau ulasan palsu.

(Input Ulasan)
Ulasan:
%s

(Konteks)
Dokumen Konteks:
%s

(Kriteria Analisis)
Pertimbangkan aspek berikut:
- Gaya penulisan ulasan
- Pengulangan frasa promosi
- Kesesuaian isi ulasan dengan produk
- Kealamian bahasa
- Ada atau tidaknya pengalaman pengguna yang spesifik
- Detail tentang kualitas produk, pengiriman, kondisi barang, warna, ukuran, bahan, atau penggunaan produk
- Pola bahasa yang terlalu umum, berlebihan, atau terlihat seperti template

Panduan klasifikasi:
- Klasifikasikan sebagai Asli jika ulasan memiliki pengalaman nyata, detail spesifik, atau informasi yang relevan dengan produk.
- Klasifikasikan sebagai Palsu jika ulasan terlalu umum, hanya berisi pujian generik, repetitif, berlebihan, atau tidak memiliki pengalaman spesifik.
- Jangan hanya mengikuti label dari dokumen konteks.
- Gunakan dokumen konteks sebagai pembanding pola bahasa dan kemiripan makna.
- Jika ulasan ambigu, pilih label yang paling kuat berdasarkan bukti pada teks dan konteks.

(Output)
Output yang dihasilkan:
Label: Asli/Palsu
Confidence: rendah/sedang/tinggi
Confidence Score: angka 0 sampai 100
Alasan: Jelaskan alasan klasifikasi secara singkat.
`, review, context)
}

func extractValue(text string, key string) string {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToLower(line), strings.ToLower(key)) {
			return strings.TrimSpace(strings.TrimPrefix(line, key))
		}
	}

	return ""
}

func extractConfidenceScore(text string) int {
	value := extractValue(text, "Confidence Score:")

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
