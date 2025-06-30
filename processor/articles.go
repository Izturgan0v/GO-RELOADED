package processor

import (
	"strings"
	"unicode"
)

// IsArticle проверяет, является ли слово артиклем.
func IsArticle(word string) bool {
	return strings.EqualFold(word, "a") || strings.EqualFold(word, "an")
}

// FixArticle корректирует артикль (a/an) на основе следующего слова.
func FixArticle(article, nextWord string) string {
	if len(nextWord) == 0 {
		return article
	}

	vowels := "aeiouhAEIOUH"
	exceptions := []string{"hour", "honest", "heir", "honor", "homage", "herb", "hysterical", "honored"}
	firstLetter := rune(nextWord[0])
	isUpper := unicode.IsUpper(rune(article[0]))

	// Проверка исключений
	for _, exception := range exceptions {
		if strings.EqualFold(nextWord, exception) {
			if isUpper {
				return "An"
			}
			return "an"
		}
	}

	// Коррекция артикля
	if strings.ContainsRune(vowels, firstLetter) && strings.EqualFold(article, "a") {
		if isUpper {
			return "An"
		}
		return "an"
	} else if !strings.ContainsRune(vowels, firstLetter) && strings.EqualFold(article, "an") {
		if isUpper {
			return "A"
		}
		return "a"
	}
	return article
}
