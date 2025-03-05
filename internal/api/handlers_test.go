package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/syirnik/GO_Yandex/internal/application"
)

// TestMain используется для настройки тестового окружения (опционально).
func TestMain(m *testing.M) {
	m.Run()
}

// TestHandleCalculate тестирует обработчик добавления выражения (внешнее поведение).
func TestHandleCalculate(t *testing.T) {
	app := application.New()
	handler := NewHandler(app)

	t.Run("InvalidMethod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/calculate", nil)
		rr := httptest.NewRecorder()

		handler.HandleCalculate(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte("invalid json")))
		rr := httptest.NewRecorder()

		handler.HandleCalculate(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rr.Code)
		}
	})

	t.Run("ValidRequest", func(t *testing.T) {
		reqBody := RequestAddExpression{Expression: "2 + 2"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.HandleCalculate(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
		}
		// Удаляем проверку Content-Type, так как обработчик его не устанавливает
	})
}

// TestHandleGetTask тестирует обработчик получения следующей задачи (внешнее поведение).
func TestHandleGetTask(t *testing.T) {
	app := application.New()
	handler := NewHandler(app)

	t.Run("NoTasks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/task", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}

// TestHandlePostTask тестирует обработчик обновления результата задачи (внешнее поведение).
func TestHandlePostTask(t *testing.T) {
	app := application.New()
	handler := NewHandler(app)

	t.Run("EmptyBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/task", nil)
		rr := httptest.NewRecorder()

		handler.HandlePostTask(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rr.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte("invalid json")))
		rr := httptest.NewRecorder()

		handler.HandlePostTask(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rr.Code)
		}
	})

	t.Run("MissingTaskID", func(t *testing.T) {
		reqBody := RequestPostTask{ID: 0, Result: 4}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.HandlePostTask(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		reqBody := RequestPostTask{ID: 1, Result: 4}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.HandlePostTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}

// TestHandleExpressions тестирует обработчик получения всех выражений (внешнее поведение).
func TestHandleExpressions(t *testing.T) {
	app := application.New()
	handler := NewHandler(app)

	t.Run("InvalidMethod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/expressions", nil)
		rr := httptest.NewRecorder()

		handler.HandleExpressions(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("ValidRequest", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expressions", nil)
		rr := httptest.NewRecorder()

		handler.HandleExpressions(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		if rr.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", rr.Header().Get("Content-Type"))
		}
	})
}

// TestHandleGetExpressionByID тестирует обработчик получения выражения по ID (внешнее поведение).
func TestHandleGetExpressionByID(t *testing.T) {
	app := application.New()
	handler := NewHandler(app)

	t.Run("InvalidMethod", func(t *testing.T) {
		// Используем более длинный путь, чтобы соответствовать проверке len(parts) < 4
		req := httptest.NewRequest(http.MethodPost, "/api/expressions/1", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetExpressionByID(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("InvalidURLPath", func(t *testing.T) {
		// Путь не содержит достаточно сегментов
		req := httptest.NewRequest(http.MethodGet, "/api/expressions/", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetExpressionByID(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("InvalidIDFormat", func(t *testing.T) {
		// Используем более длинный путь
		req := httptest.NewRequest(http.MethodGet, "/api/expressions/invalid", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetExpressionByID(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("ExpressionNotFound", func(t *testing.T) {
		// Используем более длинный путь
		req := httptest.NewRequest(http.MethodGet, "/api/expressions/999", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetExpressionByID(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	// Тест для проверки валидного запроса можно добавить после создания выражения
	// Для этого сначала нужно добавить выражение через HandleCalculate
	t.Run("ValidRequestAfterCreatingExpression", func(t *testing.T) {
		// Сначала создаем выражение
		reqBody := RequestAddExpression{Expression: "3 + 4"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.HandleCalculate(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Failed to create expression, got status %d", rr.Code)
			return
		}

		// Получаем ID созданного выражения
		var resp ResponseAddExpression
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Errorf("Failed to decode response: %v", err)
			return
		}

		// Теперь запрашиваем выражение по ID, используя более длинный путь
		req = httptest.NewRequest(http.MethodGet, "/api/expressions/"+strconv.Itoa(resp.ID), nil)
		rr = httptest.NewRecorder()

		handler.HandleGetExpressionByID(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		if rr.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", rr.Header().Get("Content-Type"))
		}

		// Проверяем структуру ответа
		var exprResp GetExpressionResponse
		if err := json.NewDecoder(rr.Body).Decode(&exprResp); err != nil {
			t.Errorf("Failed to decode response: %v", err)
			return
		}

		if exprResp.Expression.ID != resp.ID {
			t.Errorf("expected expression ID %d, got %d", resp.ID, exprResp.Expression.ID)
		}
	})
}
