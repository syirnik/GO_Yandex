package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/syirnik/GO_Yandex/pkg/calculation"
)

// Конфигурация
type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Структура Request для декодирования входящего JSON-запроса
type Request struct {
	Expression string `json:"expression"`
}

// Структура для ответа
type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Функция для проверки валидности выражения
func isValidExpression(expr string) bool {
	// Пример простого регулярного выражения для проверки
	re := regexp.MustCompile(`^[0-9+\-*/().\s]*$`)
	return re.MatchString(expr)
}

// Функция для получения сообщения об ошибке и HTTP статуса
func getErrorMessage(err error) (string, int) {
	switch err {
	case calculation.ErrEmptyExpression:
		return "Expression cannot be empty", http.StatusBadRequest
	case calculation.ErrMismatchedParentheses:
		return "Mismatched parentheses. Please check your brackets.", http.StatusUnprocessableEntity
	case calculation.ErrInsufficientOperands:
		return "Not enough operands for the operation.", http.StatusUnprocessableEntity
	case calculation.ErrDivisionByZero:
		return "Cannot divide by zero. Please check your expression.", http.StatusUnprocessableEntity
	case calculation.ErrInvalidExpression:
		return "The expression is invalid. Please check the syntax.", http.StatusUnprocessableEntity
	default:
		return "Internal server error. Please try again later.", http.StatusInternalServerError
	}
}

// Обработчик POST-запроса
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()

	// Декодируем тело запроса в структуру Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format. Please provide valid JSON."}`, http.StatusBadRequest)
		return
	}

	// Проверка на валидность выражения
	if !isValidExpression(request.Expression) {
		http.Error(w, `{"error": "Expression contains invalid characters. Only numbers, operators (+, -, *, /), and parentheses are allowed."}`, http.StatusUnprocessableEntity)
		return
	}

	// Вычисление выражения
	result, err := calculation.Calc(request.Expression)
	var response Response

	if err != nil {
		// Получаем понятное сообщение и статус для ошибки
		errorMessage, statusCode := getErrorMessage(err)
		response.Error = errorMessage
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, response.Error), statusCode)
		return
	}

	// Ответ с результатом
	response.Result = fmt.Sprintf("%f", result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Функция запуска приложения как веб-сервиса
func (a *Application) RunServer() error {
	// Регистрируем обработчик для /api/v1/calculate
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	log.Printf("Server is listening on port %s...", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
