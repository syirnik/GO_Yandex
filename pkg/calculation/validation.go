package calculation

import (
	"log"
	"unicode"
)

func ValidateExpression(expression string) error {
	log.Printf("ValidateExpression: Validating expression: %s", expression)

	// 1. Проверка на пустое выражение
	if expression == "" {
		log.Printf("ValidateExpression: Empty expression")
		return ErrEmptyExpression
	}

	// 2. Проверка на допустимые символы
	for _, char := range expression {
		if !unicode.IsDigit(char) && char != '.' && char != '+' && char != '-' &&
			char != '*' && char != '/' && char != '(' && char != ')' && !unicode.IsSpace(char) {
			log.Printf("ValidateExpression: Invalid character found: %c", char)
			return ErrInvalidCharacter
		}
	}

	// 3. Проверка баланса скобок
	openParentheses := 0
	for _, char := range expression {
		if char == '(' {
			openParentheses++
		} else if char == ')' {
			openParentheses--
			if openParentheses < 0 {
				log.Printf("ValidateExpression: Mismatched parentheses")
				return ErrMismatchedParentheses
			}
		}
	}
	if openParentheses != 0 {
		log.Printf("ValidateExpression: Mismatched parentheses")
		return ErrMismatchedParentheses
	}

	log.Printf("ValidateExpression: Expression is valid")
	return nil
}
