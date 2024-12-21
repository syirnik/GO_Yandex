package calculation

import (
	"testing"
)

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
			name:           "priority with multiplication",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	// Тесты с ожидаемым успешным результатом
	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error: %v", testCase.expression, err)
			}
			if val != testCase.expectedResult {
				t.Fatalf("test case %s: %f should be equal %f", testCase.name, val, testCase.expectedResult)
			}
		})
	}

	// Тесты с ошибками
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "invalid expression",
			expression:  "1+1*",
			expectedErr: ErrInsufficientOperands, // Ожидаем ErrInsufficientOperands
		},
		{
			name:        "mismatched parentheses (missing closing)",
			expression:  "(1+2",
			expectedErr: ErrMismatchedParentheses, // Ожидаем ErrMismatchedParentheses
		},
		{
			name:        "mismatched parentheses (missing opening)",
			expression:  "1+2)",                   // Закрывающая скобка без открывающей
			expectedErr: ErrMismatchedParentheses, // Ожидаем ErrMismatchedParentheses
		},
	}

	// Тесты с ошибками
	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Calc(testCase.expression)
			if err == nil {
				t.Fatalf("error case %s should return an error", testCase.expression)
			}
			if err != testCase.expectedErr {
				t.Fatalf("expected error %v, got %v", testCase.expectedErr, err)
			}
		})
	}

	// Тест на пустое выражение
	t.Run("empty expression", func(t *testing.T) {
		_, err := Calc("")
		if err == nil {
			t.Fatalf("empty expression should return an error")
		}
		if err != ErrEmptyExpression {
			t.Fatalf("expected error %v, got %v", ErrEmptyExpression, err)
		}
	})
}
