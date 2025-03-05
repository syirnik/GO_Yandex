package calculation

import "errors"

var (
	ErrInvalidExpression     = errors.New("invalid expression")
	ErrDivisionByZero        = errors.New("division by zero")
	ErrMismatchedParentheses = errors.New("mismatched parentheses")
	ErrInvalidCharacter      = errors.New("invalid character in expression")
	ErrEmptyExpression       = errors.New("expression is empty")
	ErrInsufficientOperands  = errors.New("insufficient operands for operation")
)
