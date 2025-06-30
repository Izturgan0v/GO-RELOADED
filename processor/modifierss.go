package processor

import (
	"strconv"
	"strings"
	"unicode"
)

// IsModifier проверяет, является ли слово модификатором.
func IsModifier(word string) bool {
	modifiers := map[string]bool{
		"cap": true,
		"up":  true,
		"low": true,
		"hex": true,
		"bin": true,
	}
	return modifiers[word] // Строгая проверка регистра
}

// IsDigit проверяет, является ли строка числом.
func IsDigit(str string) bool {
	for i, s := range str {
		if i == 0 && (s == '-' || s == '+') {
			continue
		}
		if !unicode.IsDigit(s) {
			return false
		}
	}
	return len(str) > 0 && (str != "+" && str != "-")
}

// Capitalize делает первую букву заглавной, а остальные строчными.
func Capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
}

// BaseHex конвертирует шестнадцатеричное число в десятичное.
func BaseHex(num string) string {
	decimal, err := strconv.ParseInt(num, 16, 0)
	if err != nil {
		return num
	}
	return strconv.Itoa(int(decimal))
}

// BaseBin конвертирует двоичное число в десятичное.
func BaseBin(num string) string {
	decimal, err := strconv.ParseInt(num, 2, 0)
	if err != nil {
		return num
	}
	return strconv.Itoa(int(decimal))
}

// ApplyModifier применяет модификатор к слову.
func ApplyModifier(modifier, word string) string {
	switch modifier {
	case "cap":
		return Capitalize(word)
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "hex":
		return BaseHex(word)
	case "bin":
		return BaseBin(word)
	default:
		return word
	}
}

// ParseModifier парсит модификатор и возвращает его тип, количество и флаг вложенности.
func ParseModifier(modStr string) (string, int, bool) {
	count := 1
	nested := false

	// Проверка на вложенные модификаторы
	if strings.Contains(modStr, "(") {
		nested = true
		return modStr, count, nested
	}

	// Парсинг модификатора (mod) или (mod, count)
	parts := strings.Split(modStr, ",")
	modifier := strings.TrimSpace(parts[0])

	// bin и hex с любым параметром невалидны
	if (modifier == "bin" || modifier == "hex") && len(parts) > 1 {
		return modifier, count, false
	}

	// Для cap, up, low проверяем, что параметр числовой
	if len(parts) > 1 {
		countStr := strings.TrimSpace(parts[1])
		if IsDigit(countStr) {
			count, _ = strconv.Atoi(countStr)
		} else {
			return modifier, count, false // Невалидный параметр
		}
	}
	if count < 0 {
		count = 0 // Игнорировать отрицательные значения
	}
	return modifier, count, nested
}

// ResolveNestedModifiers разрешает вложенные модификаторы, возвращая результирующий модификатор.
func ResolveNestedModifiers(modStr string) string {
	mods := strings.Split(modStr, ")")
	var modifiers []string
	for _, mod := range mods {
		mod = strings.TrimSpace(mod)
		if mod != "" {
			mod = strings.TrimPrefix(mod, "(")
			parts := strings.SplitN(mod, ",", 2)
			modifier := strings.TrimSpace(parts[0])
			if IsModifier(modifier) {
				// Пропускаем bin и hex с параметрами
				if (modifier == "bin" || modifier == "hex") && len(parts) > 1 {
					continue
				}
				count := 1
				if len(parts) > 1 {
					countStr := strings.TrimSpace(parts[1])
					if IsDigit(countStr) {
						count, _ = strconv.Atoi(countStr)
					} else {
						continue // Невалидный параметр
					}
				}
				for i := 0; i < count; i++ {
					modifiers = append(modifiers, modifier)
				}
			}
		}
	}

	result := ""
	for _, modifier := range modifiers {
		result = ApplyModifier(modifier, result)
	}
	if result == "" && len(modifiers) > 0 {
		return modifiers[len(modifiers)-1]
	}
	return result
}

// ProcessModifiers обрабатывает модификаторы в токенах, включая вложенные, и сохраняет пустые скобки.
func ProcessModifiers(tokens []string) []string {
	result := make([]string, len(tokens))
	copy(result, tokens)

	for i := 0; i < len(result); i++ {
		if result[i] == "(" {
			depth := 1
			end := i + 1
			nestedParentheses := 0
			startOfInnermost := i
			for end < len(result) && depth > 0 {
				if result[end] == "(" {
					depth++
					nestedParentheses++
					startOfInnermost = end
				} else if result[end] == ")" {
					depth--
				}
				end++
			}
			if depth != 0 {
				continue
			}
			end--

			modifierStr := strings.Join(result[startOfInnermost+1:end], "")
			if modifierStr == "" || modifierStr == "()" {
				continue
			}

			modifier, count, nested := ParseModifier(modifierStr)
			// Пропускаем, если модификатор невалиден
			if !IsModifier(modifier) && !nested {
				continue
			}
			// Пропускаем bin и hex с параметрами
			if (modifier == "bin" || modifier == "hex") && len(strings.Split(modifierStr, ",")) > 1 {
				continue
			}
			// Пропускаем cap, up, low с нечисловыми параметрами
			if (modifier == "cap" || modifier == "up" || modifier == "low") && len(strings.Split(modifierStr, ",")) > 1 {
				parts := strings.Split(modifierStr, ",")
				if len(parts) > 1 && !IsDigit(strings.TrimSpace(parts[1])) {
					continue
				}
			}

			outerParentheses := nestedParentheses - 1
			emptyParentheses := ""
			for j := 0; j < outerParentheses; j++ {
				emptyParentheses += "()"
			}

			if nested {
				modifier = ResolveNestedModifiers(modifierStr)
				for j := i - 1; count > 0 && j >= 0; j-- {
					if !IsPunctuation(result[j]) {
						isQuoted := strings.HasPrefix(result[j], "'") && strings.HasSuffix(result[j], "'")
						word := result[j]
						if isQuoted {
							word = strings.TrimPrefix(strings.TrimSuffix(word, "'"), "'")
						}
						word = ApplyModifier(modifier, word)
						if isQuoted {
							word = "'" + word + "'"
						}
						result[j] = word
						count--
					}
					if count == 0 {
						break
					}
				}
			} else {
				for j := i - 1; count > 0 && j >= 0; j-- {
					if !IsPunctuation(result[j]) {
						isQuoted := strings.HasPrefix(result[j], "'") && strings.HasSuffix(result[j], "'")
						word := result[j]
						if isQuoted {
							word = strings.TrimPrefix(strings.TrimSuffix(word, "'"), "'")
						}
						word = ApplyModifier(modifier, word)
						if isQuoted {
							word = "'" + word + "'"
						}
						result[j] = word
						count--
					}
					if count == 0 {
						break
					}
				}
			}

			result = append(result[:i], result[end+1:]...)
			if emptyParentheses != "" {
				result = append(result[:i], append([]string{emptyParentheses}, result[i:]...)...)
			}
			i--
		}
	}
	return result
}
