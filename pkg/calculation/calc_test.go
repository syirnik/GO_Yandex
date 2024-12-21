package calculation_test

import (
	"testing"

	"github.com/syirnik/GO_Yandex/pkg/calculation"
	"github.com/syirnik/GO_Yandex/pkg/calculation/errors"
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

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculation.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error: %v", testCase.expression, err)
			}
			if val != testCase.expectedResult {
				t.Fatalf("test case %s: %f should be equal %f", testCase.name, val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "simple",
			expression:  "1+1*",
			expectedErr: errors.ErrInvalidExpression, // Ожидаем ошибку
		},
		{
			name:        "divide by zero",
			expression:  "1/0",
			expectedErr: errors.ErrDivisionByZero, // Ожидаем ошибку деления на ноль
		},
		{
			name:        "incorrect parentheses",
			expression:  "(1+2",
			expectedErr: errors.ErrInvalidExpression, // Ожидаем ошибку из-за неправильных скобок
		},
		{
			name:        "incorrect parentheses",
			expression:  "1+2)",                      // Закрывающая скобка без открывающей
			expectedErr: errors.ErrInvalidExpression, // Ожидаем ошибку
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := calculation.Calc(testCase.expression)
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
		_, err := calculation.Calc("")
		if err == nil {
			t.Fatalf("empty expression should return an error")
		}
		if err.Error() != "пустое выражение" {
			t.Fatalf("expected error message: пустое выражение, got: %v", err)
		}
	})
}
