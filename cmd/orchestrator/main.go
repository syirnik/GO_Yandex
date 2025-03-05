package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/syirnik/GO_Yandex/internal/api"
	"github.com/syirnik/GO_Yandex/internal/application"

	"github.com/joho/godotenv"
)

const (
	defaultOperationTimeMs = 1000
	defaultPort            = "8080"
)

// checkEnvironmentVariable проверяет наличие и корректность переменной среды.
// Возвращает ошибку, если переменная отсутствует или некорректна.
func checkEnvironmentVariable(name string) error {
	value := os.Getenv(name)
	if value == "" {
		return fmt.Errorf("%s is not set", name)
	}
	if val, err := strconv.Atoi(value); err != nil {
		return fmt.Errorf("%s is not a valid number: %v", name, err)
	} else if val <= 0 {
		return fmt.Errorf("%s must be a positive number, got: %s", name, value)
	}
	return nil
}

// getPortFromEnv получает порт из переменной среды или возвращает дефолтное значение.
// Возвращает ошибку, если порт некорректен.
func getPortFromEnv() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Warning: PORT is not set, using default value: %s", defaultPort)
		return defaultPort, nil
	}

	// Проверяем, что порт является числом.
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("PORT is not a valid number: %v", err)
	}

	return strconv.Itoa(portNum), nil
}

func main() {
	// Загружаем переменные из .env файла, если он существует.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error: Failed to load .env file: %v", err)
	} else if os.IsNotExist(err) {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Проверяем переменные среды для времени выполнения операций.
	envVars := []string{
		"TIME_ADDITION_MS",
		"TIME_SUBTRACTION_MS",
		"TIME_MULTIPLICATIONS_MS",
		"TIME_DIVISIONS_MS",
	}
	for _, env := range envVars {
		if err := checkEnvironmentVariable(env); err != nil {
			log.Fatalf("Error: Configuration error: %v", err)
		}
	}

	// Получаем порт из переменной среды.
	port, err := getPortFromEnv()
	if err != nil {
		log.Fatalf("Error: Configuration error: %v", err)
	}

	// Создаем приложение и сервер.
	app := application.New()
	server := api.NewServer(app)

	// Запускаем сервер на указанном порте.
	log.Printf("Starting server on port :%s...", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatalf("Error: Failed to start server: %v", err)
	}

	// Логируем успешный запуск сервера.
	log.Printf("Server started on port :%s", port)
}
