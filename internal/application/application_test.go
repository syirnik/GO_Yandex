package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест успешного вычисления
func TestCalcHandler_Success(t *testing.T) {
	reqBody := `{"expression": "3 + 5 * (2 - 1)"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalcHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 but got %d", rr.Code)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	expectedResult := "8.000000" // результат вычислений с плавающей точкой

	if resp.Result != expectedResult {
		t.Errorf("expected result %s but got %s", expectedResult, resp.Result)
	}
}

// Тест деления на ноль
func TestCalcHandler_DivisionByZero(t *testing.T) {
	reqBody := `{"expression": "10 / 0"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalcHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 but got %d", rr.Code)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	expectedError := "Cannot divide by zero. Please check your expression."
	if resp.Error != expectedError {
		t.Errorf("expected error %s but got %s", expectedError, resp.Error)
	}
}

// Тест для неверного выражения (например, два оператора подряд)
func TestCalcHandler_InvalidExpression(t *testing.T) {
	reqBody := `{"expression": "3 + * 5"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalcHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 but got %d", rr.Code)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	expectedError := "Not enough operands for the operation."
	if resp.Error != expectedError {
		t.Errorf("expected error %s but got %s", expectedError, resp.Error)
	}
}

// Тест для неверного JSON
func TestCalcHandler_InvalidJSON(t *testing.T) {
	reqBody := `{"expression": "3 + 5 * (2 - 1"`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalcHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 but got %d", rr.Code)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	expectedError := "Invalid JSON format. Please provide valid JSON."
	if resp.Error != expectedError {
		t.Errorf("expected error %s but got %s", expectedError, resp.Error)
	}
}

// Тест для пустого выражения
func TestCalcHandler_EmptyExpression(t *testing.T) {
	reqBody := `{"expression": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalcHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 but got %d", rr.Code)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	expectedError := "Expression cannot be empty"
	if resp.Error != expectedError {
		t.Errorf("expected error %s but got %s", expectedError, resp.Error)
	}
}

// Тест для ошибки по умолчанию (default)
func TestCalcHandler_DefaultError(t *testing.T) {
	// Генерируем ошибку, которая не указана в switch (например, любые другие ошибки)
	err := fmt.Errorf("some unknown error")
	errorMessage, statusCode := getErrorMessage(err)

	// Проверяем, что сообщение об ошибке корректно для default
	expectedMessage := "Internal server error. Please try again later."
	if errorMessage != expectedMessage {
		t.Errorf("expected message %s but got %s", expectedMessage, errorMessage)
	}

	// Проверяем, что статус-код корректен для default
	expectedStatusCode := http.StatusInternalServerError
	if statusCode != expectedStatusCode {
		t.Errorf("expected status code %d but got %d", expectedStatusCode, statusCode)
	}
}
