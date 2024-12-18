package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// Функция для парсинга выражения — разбиваем строку на числа, операторы и скобки
func tokenize(expression string) []string {
	tokens := []string{}
	number := strings.Builder{}

	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			// Если это цифра или точка, собираем число
			number.WriteRune(char)
		} else {
			// Если встретили оператор или скобку, добавляем собранное число и сам оператор
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			if char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' {
				tokens = append(tokens, string(char))
			}
		}
	}

	// Добавляем последнее число, если оно есть
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

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
func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var operators []string

	for _, token := range tokens {
		if unicode.IsDigit(rune(token[0])) || token == "." {
			// Если это число, добавляем его в выход
			output = append(output, token)
		} else if token == "(" {
			// Открывающая скобка идёт в стек
			operators = append(operators, token)
		} else if token == ")" {
			// Выталкиваем операторы из стека до открывающей скобки
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			// Убираем открывающую скобку
			if len(operators) == 0 {
				return nil, errors.New("некорректная расстановка скобок")
			}
			operators = operators[:len(operators)-1]
		} else if token == "+" || token == "-" || token == "*" || token == "/" {
			// Работа с операторами: вытесняем из стека операторы с большим или равным приоритетом
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		}
	}

	// Переносим оставшиеся операторы в выход
	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, errors.New("некорректная расстановка скобок")
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// Функция для вычисления постфиксного выражения
func evaluatePostfix(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		if unicode.IsDigit(rune(token[0])) || token == "." {
			// Преобразуем строку в число и кладём в стек
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, value)
		} else if token == "+" || token == "-" || token == "*" || token == "/" {
			// Если это оператор, извлекаем два числа из стека
			if len(stack) < 2 {
				return 0, errors.New("некорректное выражение")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("деление на ноль")
				}
				result = a / b
			}
			stack = append(stack, result)
		}
	}

	// В конце в стеке должно остаться одно значение — результат
	if len(stack) != 1 {
		return 0, errors.New("некорректное выражение")
	}

	return stack[0], nil
}

// Основная функция калькулятора
func Calc(expression string) (float64, error) {
	// 1. Парсим выражение в токены
	tokens := tokenize(expression)

	// 2. Преобразуем инфиксное выражение в постфиксное
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}

	// 3. Вычисляем результат постфиксного выражения
	return evaluatePostfix(postfix)
}

// Функция тестирования
func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:       "simple",
			expression: "1+1*",
		},
		{
			name:       "divide by zero",
			expression: "1/0",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Calc(testCase.expression)
			if err == nil {
				t.Fatalf("error case %s should return an error", testCase.expression)
			}
		})
	}
}