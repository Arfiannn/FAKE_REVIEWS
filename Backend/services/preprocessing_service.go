package services

import (
	"regexp"
	"strings"
)

func PreprocessReviewText(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)

	urlRegex := regexp.MustCompile(`https?://\S+|www\.\S+`)
	text = urlRegex.ReplaceAllString(text, " ")

	symbolRegex := regexp.MustCompile(`[^\p{L}\p{N}\s]+`)
	text = symbolRegex.ReplaceAllString(text, " ")

	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

func IsShopeeLink(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input))

	return strings.Contains(input, "shopee.co.id") ||
		strings.Contains(input, "shp.ee")
}
