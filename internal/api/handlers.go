package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/syirnik/GO_Yandex/internal/application"
	"github.com/syirnik/GO_Yandex/pkg/calculation"
)

// Handler содержит ссылку на приложение
type Handler struct {
	App *application.Application
}

// NewHandler создает новый обработчик
func NewHandler(app *application.Application) *Handler {
	return &Handler{App: app}
}

// CORS Middleware
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// Функция для отправки JSON-ошибки клиенту
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// HandleCalculate обрабатывает добавление выражения
func (h *Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r) // Разрешаем CORS

	log.Printf("HandleCalculate: Received request")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("HandleCalculate: Invalid method %s", r.Method)
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req RequestAddExpression
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("HandleCalculate: Error decoding request body: %v", err)
		sendErrorResponse(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if err := calculation.ValidateExpression(req.Expression); err != nil {
		log.Printf("HandleCalculate: Invalid expression: %v", err)
		sendErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	exprID, err := h.App.ParseExpression(req.Expression)
	if err != nil {
		log.Printf("HandleCalculate: Error parsing expression: %v", err)
		sendErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	log.Printf("HandleCalculate: Added expression with ID: %d", exprID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponseAddExpression{ID: exprID})
}

// HandleGetTask обрабатывает запрос на получение следующей задачи агентом
func (h *Handler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	// Получаем следующую задачу из очереди
	task, err := h.App.GetNextTask()
	if err != nil {
		log.Printf("HandleGetTask: Error retrieving task: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Если задач нет, возвращаем 404
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("HandleGetTask: Returning task ID: %d", task.ID)

	// Формируем ответ
	response := ResponseGetTask{
		Task: TaskResponse{
			ID:            task.ID,
			Arg1:          task.Arg1,
			Arg2:          task.Arg2,
			Operation:     task.Operation,
			OperationTime: task.OperationTime,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("HandleGetTask: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandlePostTask принимает результат выполнения задачи от агента
func (h *Handler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandlePostTask: Received request")

	// Закрываем тело запроса при выходе из функции
	defer r.Body.Close()

	// Декодируем тело запроса
	var req RequestPostTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("HandlePostTask: Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	// Проверяем, передан ли ID задачи
	if req.ID == 0 {
		log.Printf("HandlePostTask: Missing task ID")
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	log.Printf("HandlePostTask: Processing task ID %d with result %.2f", req.ID, req.Result)

	// Обновляем результат задачи
	err := h.App.CompleteTask(req.ID, req.Result)
	if err != nil {
		log.Printf("HandlePostTask: Error completing task: %v", err)

		// Проверяем, является ли ошибка "task not found"
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("HandlePostTask: Successfully completed task ID %d", req.ID)
	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
}

// HandleExpressions обрабатывает запрос на получение всех выражений
func (h *Handler) HandleExpressions(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleExpressions: Received request")

	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		log.Printf("HandleExpressions: Invalid method %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем все выражения из приложения
	expressions := h.App.GetAllExpressions()
	log.Printf("HandleExpressions: Retrieved %d expressions", len(expressions))

	// Преобразуем выражения в JSON-массив
	var response ResponseGetExpressions
	for _, expr := range expressions {
		response.Expressions = append(response.Expressions, ExpressionResponse{
			ID:     expr.ID,
			Status: expr.Status,
			Result: expr.Result,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("HandleExpressions: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleGetExpressionByID обрабатывает GET-запрос для получения выражения по ID.
func (h *Handler) HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		log.Printf("HandleGetExpressionByID: Invalid method %s for URL: %s", r.Method, r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL пути
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[len(parts)-1] == "" {
		log.Printf("HandleGetExpressionByID: Invalid URL path: %s", r.URL.Path)
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	// Получаем ID выражения и преобразуем в число
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("HandleGetExpressionByID: Invalid ID format: %s", idStr)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Получаем выражение по ID
	expression, err := h.App.GetExpressionByID(id)
	if err != nil {
		log.Printf("HandleGetExpressionByID: Error retrieving expression ID %d: %v", id, err)

		// Проверяем, является ли ошибка "not found"
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Expression not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if expression == nil {
		log.Printf("HandleGetExpressionByID: Expression with ID %d not found", id)
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	log.Printf("HandleGetExpressionByID: Successfully retrieved expression ID %d with status %s",
		id, expression.Status)

	// Формируем JSON-ответ
	response := GetExpressionResponse{
		Expression: ExpressionResponse{
			ID:     expression.ID,
			Status: expression.Status,
			Result: expression.Result,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Проверяем ошибку при кодировании JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("HandleGetExpressionByID: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleGetResult обрабатывает запрос на получение результата выражения
func (h *Handler) HandleGetResult(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r) // Разрешаем CORS

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID выражения из URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[len(parts)-1] == "" {
		log.Printf("HandleGetResult: Invalid URL path: %s", r.URL.Path)
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	exprIDStr := parts[len(parts)-1]
	exprID, err := strconv.Atoi(exprIDStr)
	if err != nil {
		log.Printf("HandleGetResult: Invalid ID format: %s", exprIDStr)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Получаем результат выражения
	result, err := h.App.GetExpressionResult(exprID)
	if err != nil {
		log.Printf("HandleGetResult: Error retrieving result for expression ID %d: %v", exprID, err)
		http.Error(w, "Result not found", http.StatusNotFound)
		return
	}

	// Отправляем результат клиенту
	response := map[string]interface{}{
		"result": result,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
