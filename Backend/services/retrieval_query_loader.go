package services

import (
	"BE_FAKE_REVIEW/models"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadRetrievalQueries(
	csvPath string,
) ([]models.EvaluationTestReview, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuka retrieval_queries.csv: %w",
			err,
		)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membaca retrieval_queries.csv: %w",
			err,
		)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf(
			"retrieval_queries.csv kosong",
		)
	}

	headerIndexes := make(map[string]int)

	for index, column := range records[0] {
		column = strings.TrimSpace(
			strings.TrimPrefix(column, "\uFEFF"),
		)

		headerIndexes[column] = index
	}

	requiredColumns := []string{
		"query_id",
		"clean_review",
		"label",
		"label_code",
	}

	for _, column := range requiredColumns {
		if _, exists := headerIndexes[column]; !exists {
			return nil, fmt.Errorf(
				"kolom %s tidak ditemukan",
				column,
			)
		}
	}

	queryIDIndex := headerIndexes["query_id"]
	reviewIndex := headerIndexes["clean_review"]
	labelIndex := headerIndexes["label"]
	labelCodeIndex := headerIndexes["label_code"]

	queries := make(
		[]models.EvaluationTestReview,
		0,
		len(records)-1,
	)

	for rowIndex, record := range records[1:] {
		if queryIDIndex >= len(record) ||
			reviewIndex >= len(record) ||
			labelIndex >= len(record) ||
			labelCodeIndex >= len(record) {
			continue
		}

		queryID, err := strconv.ParseUint(
			strings.TrimSpace(record[queryIDIndex]),
			10,
			64,
		)
		if err != nil {
			fmt.Printf(
				"query_id tidak valid pada baris %d\n",
				rowIndex+2,
			)
			continue
		}

		labelCode, err := strconv.Atoi(
			strings.TrimSpace(record[labelCodeIndex]),
		)
		if err != nil {
			fmt.Printf(
				"label_code tidak valid pada baris %d\n",
				rowIndex+2,
			)
			continue
		}

		review := strings.TrimSpace(
			record[reviewIndex],
		)

		if review == "" {
			continue
		}

		queries = append(
			queries,
			models.EvaluationTestReview{
				ID:     uint(queryID),
				Review: review,
				Label: strings.TrimSpace(
					record[labelIndex],
				),
				LabelCode: labelCode,
			},
		)
	}

	if len(queries) == 0 {
		return nil, fmt.Errorf(
			"tidak ada query retrieval yang valid",
		)
	}

	return queries, nil
}