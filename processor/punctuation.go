package processor

import (
	"regexp"
)

// IsPunctuation проверяет, является ли токен знаком пунктуации.
func IsPunctuation(token string) bool {
	return regexp.MustCompile(`^[,.!?:;'"()]$`).MatchString(token)
}

// CorrectPunctuation корректирует пунктуацию и пробелы вокруг кавычек.
func CorrectPunctuation(text string) string {
	// Удаление пробелов перед знаками пунктуации
	re := regexp.MustCompile(`\s*([,.!?:;])`)
	text = re.ReplaceAllString(text, "$1")

	// Добавление пробела после знаков пунктуации, если он отсутствует
	text = regexp.MustCompile(`([,.!?:;])([^,.!?:;\s])`).ReplaceAllString(text, "$1 $2")

	// Коррекция пробелов внутри скобок
	reParentheses := regexp.MustCompile(`\(\s*(.*?)\s*\)`)
	text = reParentheses.ReplaceAllString(text, "($1)")

	// Коррекция двойных кавычек
	reDouble := regexp.MustCompile(`"\s*(.*?)\s*"`)
	text = reDouble.ReplaceAllString(text, "\"$1\"")

	// Коррекция одинарных кавычек: сохраняем апострофы и цитирование
	reSingle := regexp.MustCompile(`(?:\b\w*'\w+\b)|(?:'\s*([^']*?)\s*')`)
	text = reSingle.ReplaceAllStringFunc(text, func(match string) string {
		if regexp.MustCompile(`\b\w*'\w+\b`).MatchString(match) {
			return match
		}
		return regexp.MustCompile(`'\s*([^']*?)\s*'`).ReplaceAllString(match, "'$1'")
	})

	// Удаление пробелов после открывающих скобок
	reOpen := regexp.MustCompile(`\(\s*`)
	text = reOpen.ReplaceAllString(text, "(")

	// Удаление пробелов перед закрывающими скобками
	reClose := regexp.MustCompile(`\s*\)`)
	text = reClose.ReplaceAllString(text, ")")

	// Удаление лишних пробелов
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}
