package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestPerformOperation проверяет корректность выполнения математических операций.
func TestPerformOperation(t *testing.T) {
	tests := []struct {
		name      string
		arg1      float64
		arg2      float64
		operation string
		want      float64
		wantErr   bool
	}{
		// Успешные случаи
		{"Addition", 2.0, 3.0, "+", 5.0, false},
		{"Subtraction", 5.0, 3.0, "-", 2.0, false},
		{"Multiplication", 4.0, 3.0, "*", 12.0, false},
		{"Division", 6.0, 2.0, "/", 3.0, false},

		// Ошибки
		{"Division by zero", 6.0, 0.0, "/", 0.0, true},
		{"Invalid operation", 2.0, 3.0, "%", 0.0, true},
		{"NaN argument", math.NaN(), 3.0, "+", 0.0, true},
		{"Inf argument", math.Inf(1), 3.0, "+", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := performOperation(tt.arg1, tt.arg2, tt.operation)
			if (err != nil) != tt.wantErr {
				t.Errorf("performOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("performOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestWorkerSuccess проверяет успешное выполнение задачи.
func TestWorkerSuccess(t *testing.T) {
	// Создаем тестовый сервер, который возвращает задачу.
	task := Task{
		ID:            1,
		Arg1:          2.0,
		Arg2:          3.0,
		Operation:     "+",
		OperationTime: 10,
	}
	response := struct {
		Task Task `json:"task"`
	}{Task: task}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" {
			if r.Method == http.MethodGet {
				json.NewEncoder(w).Encode(response)
			} else if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
			}
		}
	}))
	defer server.Close()

	// Создаем HTTP-клиент с таймаутом.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Запускаем worker в отдельной горутине.
	done := make(chan struct{})
	go func() {
		worker(1, server.URL, client)
		close(done)
	}()

	// Ждем, чтобы worker успел обработать задачу.
	time.Sleep(100 * time.Millisecond)

	// Завершаем тест.
	server.Close()
}

// TestWorkerNoTasks проверяет обработку отсутствия задач (404).
func TestWorkerNoTasks(t *testing.T) {
	// Создаем тестовый сервер, который возвращает 404.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Создаем HTTP-клиент с таймаутом.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Запускаем worker в отдельной горутине.
	done := make(chan struct{})
	go func() {
		worker(1, server.URL, client)
		close(done)
	}()

	// Ждем, чтобы worker успел обработать отсутствие задач.
	time.Sleep(100 * time.Millisecond)

	// Завершаем тест.
	server.Close()
}
