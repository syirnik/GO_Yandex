package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Task представляет структуру задачи, полученной от оркестратора.
type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	Status        string  `json:"status"`
	Result        float64 `json:"result"`
	OperationTime int64   `json:"operation_time"`
	ParentTasks   []int   `json:"parent_tasks"`
	IsReady       bool    `json:"is_ready"`
}

// Worker представляет одну горутину, которая выполняет задачи.
func worker(workerID int, orchestratorURL string, client *http.Client) {
	log.Printf("Worker %d started and waiting for tasks...", workerID)

	// Переменная для отслеживания времени последнего лога "No tasks available".
	var lastLogTime time.Time

	for {
		// Запрашиваем задачу у оркестратора.
		resp, err := client.Get(fmt.Sprintf("%s/internal/task", orchestratorURL))
		if err != nil {
			log.Printf("Worker %d failed to fetch task from %s: %v", workerID, orchestratorURL, err)
			time.Sleep(5 * time.Second) // Ждем перед повторной попыткой.
			continue
		}

		// Проверяем статус ответа.
		if resp.StatusCode == http.StatusNotFound {
			// Логируем "No tasks available" только раз в минуту.
			if time.Since(lastLogTime) > time.Minute {
				log.Printf("Worker %d: No tasks available from %s. Waiting...", workerID, orchestratorURL)
				lastLogTime = time.Now()
			}
			resp.Body.Close()
			time.Sleep(5 * time.Second) // Ждем перед следующим запросом.
			continue
		}

		// Проверяем неожиданные коды состояния.
		if resp.StatusCode != http.StatusOK {
			log.Printf("Worker %d received unexpected status code %d from %s", workerID, resp.StatusCode, orchestratorURL)
			resp.Body.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Декодируем задачу из ответа.
		var response struct {
			Task Task `json:"task"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Printf("Worker %d failed to decode task from %s: %v", workerID, orchestratorURL, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close() // Закрываем тело ответа, чтобы избежать утечек.

		task := response.Task
		log.Printf("Worker %d received task: ID=%d, Operation=%s, Arg1=%.2f, Arg2=%.2f, OperationTime=%dms",
			workerID, task.ID, task.Operation, task.Arg1, task.Arg2, task.OperationTime)

		// Выполняем задачу с учетом задержки.
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result, err := performOperation(task.Arg1, task.Arg2, task.Operation)
		if err != nil {
			log.Printf("Worker %d failed to perform task %d: %v", workerID, task.ID, err)
			continue
		}

		// Подготавливаем результат для отправки.
		responseData := map[string]interface{}{
			"id":     task.ID,
			"result": result,
		}
		jsonBody, err := json.Marshal(responseData)
		if err != nil {
			log.Printf("Worker %d failed to marshal response for task %d: %v", workerID, task.ID, err)
			continue
		}

		// Отправляем результат обратно оркестратору.
		resp, err = client.Post(fmt.Sprintf("%s/internal/task", orchestratorURL), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Printf("Worker %d failed to send result for task %d to %s: %v", workerID, task.ID, orchestratorURL, err)
			continue
		}
		resp.Body.Close()

		log.Printf("Worker %d successfully sent result for task %d to %s", workerID, task.ID, orchestratorURL)
	}
}

// performOperation выполняет математическую операцию над аргументами.
func performOperation(arg1, arg2 float64, operation string) (float64, error) {
	// Проверяем на NaN и бесконечные значения.
	if math.IsNaN(arg1) || math.IsNaN(arg2) || math.IsInf(arg1, 0) || math.IsInf(arg2, 0) {
		return 0, fmt.Errorf("invalid arguments: NaN or Inf")
	}

	switch operation {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return arg1 / arg2, nil
	default:
		return 0, fmt.Errorf("invalid operation: %s", operation)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Получаем PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	// Получаем ORCHESTRATOR_URL и подставляем PORT
	orchestratorURLTemplate := os.Getenv("ORCHESTRATOR_URL")
	if orchestratorURLTemplate == "" {
		log.Fatal("ORCHESTRATOR_URL environment variable is not set")
	}
	orchestratorURL := fmt.Sprintf(orchestratorURLTemplate, port)
	log.Println("ORCHESTRATOR_URL:", orchestratorURL) // Логируем итоговый URL

	if _, err := url.Parse(orchestratorURL); err != nil {
		log.Fatalf("Invalid ORCHESTRATOR_URL: %v", err)
	}

	// Создаем HTTP-клиент с таймаутом.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Получаем количество горутин из переменной среды.
	computingPower := 1
	if computingPowerStr := os.Getenv("COMPUTING_POWER"); computingPowerStr != "" {
		var err error
		computingPower, err = strconv.Atoi(computingPowerStr)
		if err != nil {
			log.Fatalf("Invalid COMPUTING_POWER value: %v", err)
		}
	}
	if computingPower <= 0 {
		log.Fatalf("COMPUTING_POWER must be positive, got: %d", computingPower)
	}

	log.Printf("Agent started with computing power: %d", computingPower)

	// Запускаем пул горутин.
	for i := 0; i < computingPower; i++ {
		go worker(i+1, orchestratorURL, client)
	}

	// Блокируем основной поток, чтобы программа продолжала работать.
	select {}
}
