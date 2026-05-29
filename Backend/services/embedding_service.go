package services

import (
	"BE_FAKE_REVIEW/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EmbeddingService struct {
	Config config.Config
}

type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func NewEmbeddingService(cfg config.Config) EmbeddingService {
	return EmbeddingService{
		Config: cfg,
	}
}

func (s EmbeddingService) CreateEmbedding(text string) ([]float32, error) {
	body := OllamaEmbeddingRequest{
		Model:  s.Config.OllamaEmbedModel,
		Prompt: text,
	}

	jsonBody, _ := json.Marshal(body)

	url := fmt.Sprintf("%s/api/embeddings", s.Config.OllamaURL)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gagal membuat embedding: %s", string(responseBody))
	}

	var result OllamaEmbeddingResponse

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, err
	}

	if len(result.Embedding) == 0 {
		return nil, errors.New("embedding kosong")
	}

	if len(result.Embedding) != 768 {
		return nil, fmt.Errorf("dimensi embedding tidak sesuai, dapat %d, harus 768", len(result.Embedding))
	}

	return result.Embedding, nil
}
