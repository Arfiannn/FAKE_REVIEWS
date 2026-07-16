package services

import (
	"BE_FAKE_REVIEW/models"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadEvaluationTestData(
	csvPath string,
) ([]models.EvaluationTestReview, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuka dataset test %s: %w",
			csvPath,
			err,
		)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membaca dataset test: %w",
			err,
		)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf(
			"dataset test kosong atau tidak memiliki data",
		)
	}

	header := records[0]

	cleanReviewIndex := findCSVColumnIndex(
		header,
		"clean_review",
	)

	labelIndex := findCSVColumnIndex(
		header,
		"label",
	)

	labelCodeIndex := findCSVColumnIndex(
		header,
		"label_code",
	)

	if cleanReviewIndex == -1 {
		return nil, fmt.Errorf(
			"kolom clean_review tidak ditemukan",
		)
	}

	if labelIndex == -1 {
		return nil, fmt.Errorf(
			"kolom label tidak ditemukan",
		)
	}

	if labelCodeIndex == -1 {
		return nil, fmt.Errorf(
			"kolom label_code tidak ditemukan",
		)
	}

	testReviews := make(
		[]models.EvaluationTestReview,
		0,
	)

	for rowIndex, record := range records[1:] {
		maxIndex := cleanReviewIndex

		if labelIndex > maxIndex {
			maxIndex = labelIndex
		}

		if labelCodeIndex > maxIndex {
			maxIndex = labelCodeIndex
		}

		if len(record) <= maxIndex {
			fmt.Printf(
				"Baris %d dilewati karena kolom tidak lengkap\n",
				rowIndex+2,
			)
			continue
		}

		cleanReview := strings.TrimSpace(
			record[cleanReviewIndex],
		)

		label := strings.TrimSpace(
			record[labelIndex],
		)

		labelCodeText := strings.TrimSpace(
			record[labelCodeIndex],
		)

		if cleanReview == "" {
			fmt.Printf(
				"Baris %d dilewati karena clean_review kosong\n",
				rowIndex+2,
			)
			continue
		}

		labelCode, err := strconv.Atoi(labelCodeText)
		if err != nil {
			fmt.Printf(
				"Baris %d dilewati karena label_code tidak valid\n",
				rowIndex+2,
			)
			continue
		}

		if labelCode != 0 && labelCode != 1 {
			fmt.Printf(
				"Baris %d dilewati karena label_code bukan 0 atau 1\n",
				rowIndex+2,
			)
			continue
		}

		if label == "" {
			if labelCode == 1 {
				label = "Asli"
			} else {
				label = "Palsu"
			}
		}

		testReviews = append(
			testReviews,
			models.EvaluationTestReview{
				ID:        uint(len(testReviews) + 1),
				Review:    cleanReview,
				Label:     label,
				LabelCode: labelCode,
			},
		)
	}

	if len(testReviews) == 0 {
		return nil, fmt.Errorf(
			"tidak ada data test yang valid",
		)
	}

	return testReviews, nil
}

func findCSVColumnIndex(
	header []string,
	columnName string,
) int {
	for index, column := range header {
		cleanColumn := strings.TrimSpace(
			strings.TrimPrefix(
				column,
				"\uFEFF",
			),
		)

		if strings.EqualFold(
			cleanColumn,
			columnName,
		) {
			return index
		}
	}

	return -1
}
