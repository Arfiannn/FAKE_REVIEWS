package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func GenerateRetrievalGroundTruth(
	annotationPath string,
	outputPath string,
) (int, int, error) {
	file, err := os.Open(annotationPath)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"gagal membuka file anotasi: %w",
			err,
		)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return 0, 0, fmt.Errorf(
			"gagal membaca file anotasi: %w",
			err,
		)
	}

	if len(records) < 2 {
		return 0, 0, fmt.Errorf(
			"file anotasi kosong",
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
		"document_id",
		"is_relevant",
	}

	for _, column := range requiredColumns {
		if _, exists := headerIndexes[column]; !exists {
			return 0, 0, fmt.Errorf(
				"kolom %s tidak ditemukan",
				column,
			)
		}
	}

	queryIDIndex := headerIndexes["query_id"]
	documentIDIndex := headerIndexes["document_id"]
	isRelevantIndex := headerIndexes["is_relevant"]

	groundTruth := make(
		map[uint]map[uint]struct{},
	)

	totalRelevantDocuments := 0

	for _, record := range records[1:] {
		if queryIDIndex >= len(record) ||
			documentIDIndex >= len(record) ||
			isRelevantIndex >= len(record) {
			continue
		}

		isRelevant := strings.TrimSpace(
			record[isRelevantIndex],
		)

		if isRelevant != "1" {
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

		documentID, err := strconv.ParseUint(
			strings.TrimSpace(record[documentIDIndex]),
			10,
			64,
		)
		if err != nil {
			continue
		}

		queryKey := uint(queryID)
		documentKey := uint(documentID)

		if _, exists := groundTruth[queryKey]; !exists {
			groundTruth[queryKey] =
				make(map[uint]struct{})
		}

		if _, exists :=
			groundTruth[queryKey][documentKey]; exists {
			continue
		}

		groundTruth[queryKey][documentKey] =
			struct{}{}

		totalRelevantDocuments++
	}

	if len(groundTruth) == 0 {
		return 0, 0, fmt.Errorf(
			"tidak ada dokumen yang ditandai relevan",
		)
	}

	outputDirectory := filepath.Dir(outputPath)

	if outputDirectory != "." {
		if err := os.MkdirAll(
			outputDirectory,
			os.ModePerm,
		); err != nil {
			return 0, 0, err
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"gagal membuat ground truth: %w",
			err,
		)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write([]string{
		"query_id",
		"relevant_document_ids",
	}); err != nil {
		return 0, 0, err
	}

	queryIDs := make([]int, 0, len(groundTruth))

	for queryID := range groundTruth {
		queryIDs = append(
			queryIDs,
			int(queryID),
		)
	}

	sort.Ints(queryIDs)

	for _, queryID := range queryIDs {
		documentIDs := make(
			[]int,
			0,
			len(groundTruth[uint(queryID)]),
		)

		for documentID :=
			range groundTruth[uint(queryID)] {
			documentIDs = append(
				documentIDs,
				int(documentID),
			)
		}

		sort.Ints(documentIDs)

		documentIDTexts := make(
			[]string,
			0,
			len(documentIDs),
		)

		for _, documentID := range documentIDs {
			documentIDTexts = append(
				documentIDTexts,
				strconv.Itoa(documentID),
			)
		}

		if err := writer.Write([]string{
			strconv.Itoa(queryID),
			strings.Join(documentIDTexts, ","),
		}); err != nil {
			return 0, 0, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return 0, 0, err
	}

	return len(groundTruth),
		totalRelevantDocuments,
		nil
}