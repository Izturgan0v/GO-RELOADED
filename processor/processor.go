package processor

import (
	"strings"
)

// ProcessLine обрабатывает входную строку, выполняя токенизацию, обработку модификаторов и коррекцию артиклей и пунктуации.
func ProcessLine(line string) string {
	// Токенизация
	tokens := Split(line)
	// Обработка модификаторов
	tokens = ProcessModifiers(tokens)
	// Построение строки с коррекцией артиклей
	var text strings.Builder
	for i, token := range tokens {
		if i < len(tokens)-1 && IsArticle(token) {
			token = FixArticle(token, tokens[i+1])
		}
		if token != "" {
			text.WriteString(token + " ")
		}
	}
	// Коррекция пунктуации и кавычек
	return strings.TrimSpace(CorrectPunctuation(text.String()))
}
