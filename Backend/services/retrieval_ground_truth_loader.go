package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadRetrievalGroundTruth(
	csvPath string,
) (map[uint][]uint, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuka ground truth: %w",
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
			"gagal membaca ground truth: %w",
			err,
		)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf(
			"ground truth kosong",
		)
	}

	headerIndexes := make(map[string]int)

	for index, column := range records[0] {
		column = strings.TrimSpace(
			strings.TrimPrefix(column, "\uFEFF"),
		)

		headerIndexes[column] = index
	}

	queryIDIndex, queryExists :=
		headerIndexes["query_id"]

	documentIDsIndex, documentExists :=
		headerIndexes["relevant_document_ids"]

	if !queryExists || !documentExists {
		return nil, fmt.Errorf(
			"header ground truth tidak valid",
		)
	}

	groundTruth := make(map[uint][]uint)

	for _, record := range records[1:] {
		if queryIDIndex >= len(record) ||
			documentIDsIndex >= len(record) {
			continue
		}

		queryID, err := strconv.ParseUint(
			strings.TrimSpace(record[queryIDIndex]),
			10,
			64,
		)
		if err != nil {
			continue
		}

		documentIDTexts := strings.Split(
			record[documentIDsIndex],
			",",
		)

		documentIDs := make([]uint, 0)

		for _, documentIDText := range documentIDTexts {
			documentID, err := strconv.ParseUint(
				strings.TrimSpace(documentIDText),
				10,
				64,
			)
			if err != nil {
				continue
			}

			documentIDs = append(
				documentIDs,
				uint(documentID),
			)
		}

		if len(documentIDs) > 0 {
			groundTruth[uint(queryID)] =
				documentIDs
		}
	}

	if len(groundTruth) == 0 {
		return nil, fmt.Errorf(
			"tidak ada ground truth valid",
		)
	}

	return groundTruth, nil
}

func createRelevantDocumentSet(
	documentIDs []uint,
) map[uint]struct{} {
	result := make(map[uint]struct{})

	for _, documentID := range documentIDs {
		result[documentID] = struct{}{}
	}

	return result
}
