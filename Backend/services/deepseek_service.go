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

type DeepSeekService struct {
	Config config.Config
}

type DeepSeekChatRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepSeekMessage `json:"messages"`
	Temperature float64           `json:"temperature"`
}

type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekChatResponse struct {
	Choices []struct {
		Message DeepSeekMessage `json:"message"`
	} `json:"choices"`
}

func NewDeepSeekService(cfg config.Config) DeepSeekService {
	return DeepSeekService{
		Config: cfg,
	}
}

func (s DeepSeekService) Chat(prompt string) (string, error) {
	if s.Config.DeepSeekAPIKey == "" {
		return "", errors.New("DEEPSEEK_API_KEY belum diisi di .env")
	}

	body := DeepSeekChatRequest{
		Model: s.Config.DeepSeekModel,
		Messages: []DeepSeekMessage{
			{
				Role:    "system",
				Content: "Kamu adalah AI untuk analisis fake review e-commerce berbahasa Indonesia.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.deepseek.com/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Config.DeepSeekAPIKey)

	client := http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("gagal call DeepSeek: %s", string(responseBody))
	}

	var result DeepSeekChatResponse

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("response DeepSeek kosong")
	}

	return result.Choices[0].Message.Content, nil
}
