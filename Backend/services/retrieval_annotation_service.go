package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func (s EvaluationService) GenerateRetrievalAnnotation(
	queryCSVPath string,
	outputCSVPath string,
	limit int,
	candidateK int,
) (int, int, error) {
	if limit <= 0 {
		limit = 100
	}

	if candidateK <= 0 {
		candidateK = 20
	}

	queries, err := LoadRetrievalQueries(
		queryCSVPath,
	)
	if err != nil {
		return 0, 0, err
	}

	if limit < len(queries) {
		queries = queries[:limit]
	}

	outputDirectory := filepath.Dir(outputCSVPath)

	if outputDirectory != "." {
		if err := os.MkdirAll(
			outputDirectory,
			os.ModePerm,
		); err != nil {
			return 0, 0, fmt.Errorf(
				"gagal membuat folder output: %w",
				err,
			)
		}
	}

	file, err := os.Create(outputCSVPath)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"gagal membuat retrieval_annotation.csv: %w",
			err,
		)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"query_id",
		"query_review",
		"query_label",
		"query_label_code",
		"document_id",
		"document_review",
		"document_label",
		"document_label_code",
		"similarity_score",
		"is_relevant",
	}

	if err := writer.Write(header); err != nil {
		return 0, 0, err
	}

	processedQueries := 0
	totalRows := 0

	for index, query := range queries {
		fmt.Printf(
			"Membuat kandidat %d/%d | query_id=%d\n",
			index+1,
			len(queries),
			query.ID,
		)

		results, err := s.RAGService.VectorDB.
			SearchSimilarByText(
				query.Review,
				candidateK,
			)

		if err != nil {
			fmt.Printf(
				"Gagal query_id=%d: %v\n",
				query.ID,
				err,
			)
			continue
		}

		if len(results) == 0 {
			continue
		}

		processedQueries++

		for _, result := range results {
			row := []string{
				strconv.FormatUint(
					uint64(query.ID),
					10,
				),
				query.Review,
				query.Label,
				strconv.Itoa(query.LabelCode),

				// Menggunakan reviews.id sebagai document_id.
				strconv.FormatUint(
					uint64(result.ID),
					10,
				),

				result.CleanReview,
				result.Label,
				strconv.Itoa(result.LabelCode),

				strconv.FormatFloat(
					result.Similarity,
					'f',
					8,
					64,
				),

				// Diisi setelah file berhasil dibuat.
				"",
			}

			if err := writer.Write(row); err != nil {
				return processedQueries, totalRows, err
			}

			totalRows++
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return processedQueries, totalRows, err
	}

	return processedQueries, totalRows, nil
}
