package calculation

import (
	"log"
	"strconv"
	"strings"
	"unicode"
)

// Функция для парсинга выражения — разбиваем строку на числа, операторы и скобки
func Tokenize(expression string) []string {
	log.Printf("Tokenize: Received expression: %s", expression)

	tokens := []string{}
	number := strings.Builder{}
	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			number.WriteRune(char)
		} else {
			if number.Len() > 0 {
				token := number.String()
				log.Printf("Tokenize: Found number: %s", token)
				tokens = append(tokens, token)
				number.Reset()
			}
			if char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' {
				token := string(char)
				log.Printf("Tokenize: Found operator/parenthesis: %s", token)
				tokens = append(tokens, token)
			}
		}
	}
	if number.Len() > 0 {
		token := number.String()
		log.Printf("Tokenize: Found final number: %s", token)
		tokens = append(tokens, token)
	}

	log.Printf("Tokenize: Final tokens: %v", tokens)
	return tokens
}

// Функция для определения приоритета операторов
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// Алгоритм сортировочной станции — преобразование инфиксного выражения в постфиксное
func InfixToPostfix(tokens []string) ([]string, error) {
	log.Printf("InfixToPostfix: Converting tokens to postfix: %v", tokens)

	var output []string
	var operators []string
	for _, token := range tokens {
		if unicode.IsDigit(rune(token[0])) || token == "." {
			// Если это число, добавляем его в выход
			log.Printf("InfixToPostfix: Adding number to output: %s", token)
			output = append(output, token)
		} else if token == "(" {
			// Открывающая скобка идёт в стек
			log.Printf("InfixToPostfix: Pushing opening parenthesis to stack: %s", token)
			operators = append(operators, token)
		} else if token == ")" {
			// Выталкиваем операторы из стека до открывающей скобки
			log.Printf("InfixToPostfix: Processing closing parenthesis: %s", token)
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				op := operators[len(operators)-1]
				log.Printf("InfixToPostfix: Popping operator from stack: %s", op)
				output = append(output, op)
				operators = operators[:len(operators)-1]
			}
			// Убираем открывающую скобку
			if len(operators) == 0 {
				log.Printf("InfixToPostfix: Error: Mismatched parentheses")
				return nil, ErrMismatchedParentheses
			}
			log.Printf("InfixToPostfix: Removing opening parenthesis from stack")
			operators = operators[:len(operators)-1]
		} else if token == "+" || token == "-" || token == "*" || token == "/" {
			// Работа с операторами: вытесняем из стека операторы с большим или равным приоритетом
			log.Printf("InfixToPostfix: Processing operator: %s", token)
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				op := operators[len(operators)-1]
				log.Printf("InfixToPostfix: Popping higher precedence operator from stack: %s", op)
				output = append(output, op)
				operators = operators[:len(operators)-1]
			}
			log.Printf("InfixToPostfix: Pushing operator to stack: %s", token)
			operators = append(operators, token)
		}
	}
	// Переносим оставшиеся операторы в выход
	log.Printf("InfixToPostfix: Transferring remaining operators from stack to output")
	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			log.Printf("InfixToPostfix: Error: Mismatched parentheses")
			return nil, ErrMismatchedParentheses
		}
		op := operators[len(operators)-1]
		log.Printf("InfixToPostfix: Popping remaining operator from stack: %s", op)
		output = append(output, op)
		operators = operators[:len(operators)-1]
	}

	log.Printf("InfixToPostfix: Final postfix expression: %v", output)
	return output, nil
}

// Функция для вычисления постфиксного выражения
// Функция для вычисления постфиксного выражения
func EvaluatePostfix(tokens []string) (float64, error) {
	log.Printf("EvaluatePostfix: Evaluating postfix expression: %v", tokens)

	var stack []float64
	for _, token := range tokens {
		if unicode.IsDigit(rune(token[0])) || token == "." {
			// Преобразуем строку в число и кладём в стек
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				log.Printf("EvaluatePostfix: Error parsing number: %v", err)
				return 0, ErrInvalidCharacter // Используем ErrInvalidCharacter для недопустимых чисел
			}
			log.Printf("EvaluatePostfix: Pushing number to stack: %.2f", value)
			stack = append(stack, value)
		} else if token == "+" || token == "-" || token == "*" || token == "/" {
			// Если это оператор, извлекаем два числа из стека
			if len(stack) < 2 {
				log.Printf("EvaluatePostfix: Error: Insufficient operands for operator: %s", token)
				return 0, ErrInsufficientOperands // Используем ErrInsufficientOperands
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			// Выполняем операцию
			var result float64
			switch token {
			case "+":
				result = a + b
				log.Printf("EvaluatePostfix: Performing addition: %.2f + %.2f = %.2f", a, b, result)
			case "-":
				result = a - b
				log.Printf("EvaluatePostfix: Performing subtraction: %.2f - %.2f = %.2f", a, b, result)
			case "*":
				result = a * b
				log.Printf("EvaluatePostfix: Performing multiplication: %.2f * %.2f = %.2f", a, b, result)
			case "/":
				if b == 0 {
					log.Printf("EvaluatePostfix: Error: Division by zero")
					return 0, ErrDivisionByZero // Используем ErrDivisionByZero
				}
				result = a / b
				log.Printf("EvaluatePostfix: Performing division: %.2f / %.2f = %.2f", a, b, result)
			}

			// Добавляем результат обратно в стек
			log.Printf("EvaluatePostfix: Pushing result to stack: %.2f", result)
			stack = append(stack, result)
		} else {
			log.Printf("EvaluatePostfix: Error: Invalid token: %s", token)
			return 0, ErrInvalidCharacter // Используем ErrInvalidCharacter для недопустимых токенов
		}
	}

	// В конце в стеке должно остаться одно значение — результат
	if len(stack) != 1 {
		log.Printf("EvaluatePostfix: Error: Invalid expression format")
		return 0, ErrInvalidExpression // Используем ErrInvalidExpression
	}

	result := stack[0]
	log.Printf("EvaluatePostfix: Final result: %.2f", result)
	return result, nil
}

// Основная функция калькулятора
func Calc(expression string) (float64, error) {
	log.Printf("Calc: Starting calculation for expression: %s", expression)

	// Проверка на пустое выражение
	if expression == "" {
		log.Printf("Calc: Error: Empty expression")
		return 0, ErrEmptyExpression
	}

	// 1. Парсим выражение в токены
	tokens := Tokenize(expression)

	// 2. Преобразуем инфиксное выражение в постфиксное
	postfix, err := InfixToPostfix(tokens)
	if err != nil {
		log.Printf("Calc: Error converting to postfix: %v", err)
		return 0, err
	}

	// 3. Вычисляем результат постфиксного выражения
	result, err := EvaluatePostfix(postfix)
	if err != nil {
		log.Printf("Calc: Error evaluating postfix: %v", err)
		return 0, err
	}

	log.Printf("Calc: Final result: %.2f", result)
	return result, nil
}
