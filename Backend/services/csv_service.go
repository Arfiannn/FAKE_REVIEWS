package services

import (
	"BE_FAKE_REVIEW/models"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

func LoadReviewsFromCSV(path string) ([]models.Review, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headerMap := make(map[string]int)

	for i, h := range headers {
		headerMap[strings.TrimSpace(h)] = i
	}

	var reviews []models.Review
	seen := map[string]bool{}

	for {
		record, err := reader.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			continue
		}

		cleanReview := getCSVValue(record, headerMap, "clean_review")
		if cleanReview == "" {
			continue
		}

		if seen[cleanReview] {
			continue
		}

		seen[cleanReview] = true

		rating, _ := strconv.Atoi(getCSVValue(record, headerMap, "rating"))
		labelCode, _ := strconv.Atoi(getCSVValue(record, headerMap, "label_code"))

		review := models.Review{
			ProductName: getCSVValue(record, headerMap, "product_name"),
			ShopName:    getCSVValue(record, headerMap, "shop_name"),
			Username:    getCSVValue(record, headerMap, "username"),
			Rating:      rating,
			Review:      getCSVValue(record, headerMap, "review"),
			CleanReview: cleanReview,
			Label:       getCSVValue(record, headerMap, "label"),
			LabelCode:   labelCode,
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

func getCSVValue(record []string, headerMap map[string]int, key string) string {
	index, ok := headerMap[key]
	if !ok || index >= len(record) {
		return ""
	}

	return strings.TrimSpace(record[index])
}
