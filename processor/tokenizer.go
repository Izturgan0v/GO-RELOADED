package processor

import (
	"strings"
	"unicode"
)

// Split разбивает строку на токены, сохраняя пунктуацию.
func Split(line string) []string {
	var result []string
	var token strings.Builder

	for _, char := range line {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '+' || char == '-' {
			token.WriteRune(char)
		} else {
			if token.Len() > 0 {
				result = append(result, token.String())
				token.Reset()
			}
			if !unicode.IsSpace(char) {
				result = append(result, string(char))
			}
		}
	}

	if token.Len() > 0 {
		result = append(result, token.String())
	}
	return result
}
