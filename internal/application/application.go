package application

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/syirnik/GO_Yandex/pkg/calculation"
)

// Task представляет одну задачу (операцию)
type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	Status        string  `json:"status"`
	Result        float64 `json:"result"`
	OperationTime int64   `json:"operation_time"` // Время выполнения задачи
	ParentTasks   []int   `json:"parent_tasks"`   // ID родительских задач
	IsReady       bool    `json:"is_ready"`       // Флаг готовности задачи
}

// Expression представляет выражение, состоящее из задач
type Expression struct {
	ID     int
	Value  string
	Tasks  []*Task
	Status string
	Result float64
}

// Application управляет очередью задач и выражениями
type Application struct {
	nextTaskID       int
	nextExpressionID int
	expressions      map[int]*Expression
	taskToExpression map[int]int
	taskQueue        []*Task         // Очередь готовых задач
	dependentQueue   []*Task         // Очередь зависимых задач
	taskResults      map[int]float64 // Хранение выполненных задач
	mu               sync.Mutex
}

// New создает новый экземпляр Application
func New() *Application {
	return &Application{
		nextTaskID:       1,
		nextExpressionID: 1,
		expressions:      make(map[int]*Expression),
		taskToExpression: make(map[int]int),
		taskQueue:        []*Task{},
		dependentQueue:   []*Task{},
		taskResults:      make(map[int]float64),
	}
}

// ParseExpression разбирает выражение и создает задачи
func (app *Application) ParseExpression(expression string) (int, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	log.Printf("ParseExpression: Processing expression: %s", expression)

	tokens := calculation.Tokenize(expression)
	log.Printf("ParseExpression: Tokenized expression: %v", tokens)

	postfix, err := calculation.InfixToPostfix(tokens)
	if err != nil {
		return 0, fmt.Errorf("error converting to postfix: %w", err)
	}

	log.Printf("ParseExpression: Postfix notation: %v", postfix)

	var tasks []*Task
	stack := []int{}

	// Для каждой задачи в постфиксной записи
	for _, token := range postfix {
		log.Printf("ParseExpression: Processing token: %s", token)

		if isNumber(token) {
			// Число — это уже "выполненная задача"
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number format in expression: %s", token)
			}
			taskID := app.nextTaskID
			app.taskResults[taskID] = value // Записываем число как результат
			stack = append(stack, taskID)
			log.Printf("ParseExpression: Created task for number %s with ID %d", token, taskID)
			app.nextTaskID++
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression format: not enough values in stack")
			}

			// Извлекаем операнды для текущей операции
			task2ID := stack[len(stack)-1] // Второй операнд
			task1ID := stack[len(stack)-2] // Первый операнд
			stack = stack[:len(stack)-2]

			log.Printf("ParseExpression: Creating task for operator %s, using task IDs: %d, %d", token, task1ID, task2ID)

			// Создание задачи на операцию
			task := &Task{
				ID:            app.nextTaskID,
				Operation:     token,
				Status:        "pending",
				ParentTasks:   []int{task1ID, task2ID},
				OperationTime: app.getOperationTime(token),
			}

			// Присваиваем аргументы, если они уже вычислены
			if res1, exists1 := app.taskResults[task1ID]; exists1 {
				task.Arg1 = res1
				log.Printf("ParseExpression: Arg1 for task ID %d set to %f", task.ID, res1)
			}
			if res2, exists2 := app.taskResults[task2ID]; exists2 {
				task.Arg2 = res2
				log.Printf("ParseExpression: Arg2 for task ID %d set to %f", task.ID, res2)
			}
			// Проверяем, известны ли оба аргумента
			_, exists1 := app.taskResults[task1ID]
			_, exists2 := app.taskResults[task2ID]
			if exists1 && exists2 {
				// Проверка деления на ноль
				if token == "/" && task.Arg2 == 0 {
					return 0, errors.New("division by zero detected")
				}
				task.IsReady = true
				log.Printf("ParseExpression: Task ID %d is ready (IsReady = %v)", task.ID, task.IsReady)
				app.taskQueue = append(app.taskQueue, task)
				log.Printf("ParseExpression: Task ID %d added to taskQueue", task.ID)
			} else {
				app.dependentQueue = append(app.dependentQueue, task) // Иначе откладываем задачу
				log.Printf("ParseExpression: Task ID %d added to dependentQueue", task.ID)
			}

			tasks = append(tasks, task)
			stack = append(stack, task.ID)
			log.Printf("ParseExpression: Task ID %d added to stack", task.ID)
			app.nextTaskID++
		}
	}

	// Если все задачи созданы правильно
	if len(stack) != 1 {
		return 0, errors.New("invalid expression: stack not reduced to single result")
	}

	// Создаем выражение
	exprID := app.nextExpressionID
	app.nextExpressionID++

	expr := &Expression{
		ID:     exprID,
		Value:  expression,
		Tasks:  tasks,
		Status: "pending",
	}

	app.expressions[exprID] = expr
	log.Printf("ParseExpression: Created expression ID %d with %d tasks", exprID, len(tasks))
	return exprID, nil
}

// CompleteTask принимает результат выполнения задачи
func (app *Application) CompleteTask(taskID int, result float64) error {
	app.mu.Lock()
	defer app.mu.Unlock()

	// Сохраняем результат в карте результатов
	app.taskResults[taskID] = result
	log.Printf("CompleteTask: Task ID %d completed with result: %.2f", taskID, result)

	// Находим задачу в выражении и обновляем её статус и результат
	var taskUpdated bool
	for exprID, expr := range app.expressions {
		for _, task := range expr.Tasks {
			if task.ID == taskID {
				task.Status = "completed"
				task.Result = result
				taskUpdated = true
				log.Printf("CompleteTask: Updated task ID %d status to 'completed' in expression ID %d", taskID, exprID)

				// Проверяем, все ли задачи в выражении завершены
				allCompleted := true
				for _, t := range expr.Tasks {
					if t.Status != "completed" {
						allCompleted = false
						break
					}
				}

				// Если все задачи завершены, обновляем статус выражения
				if allCompleted {
					// Последняя задача в списке должна содержать финальный результат
					finalTask := expr.Tasks[len(expr.Tasks)-1]
					expr.Status = "completed"
					expr.Result = finalTask.Result
					log.Printf("CompleteTask: Expression ID %d completed with result %.2f", exprID, expr.Result)
				}
				break
			}
		}
		if taskUpdated {
			break
		}
	}

	// Проверяем, была ли найдена и обновлена задача
	if !taskUpdated {
		log.Printf("CompleteTask: Task ID %d not found", taskID)
		return fmt.Errorf("task with ID %d not found", taskID)
	}

	// Обновляем зависимые задачи
	var newDependentQueue []*Task
	for _, task := range app.dependentQueue {
		ready := true
		log.Printf("CompleteTask: Checking dependencies for task ID %d", task.ID)

		// Проверяем все родительские задачи
		for i, parentID := range task.ParentTasks {
			if res, exists := app.taskResults[parentID]; exists {
				// Подставляем результаты в соответствующие аргументы задачи
				if i == 0 {
					task.Arg1 = res
					log.Printf("CompleteTask: Arg1 for task ID %d set to %f from parent ID %d", task.ID, res, parentID)
				} else if i == 1 {
					task.Arg2 = res
					log.Printf("CompleteTask: Arg2 for task ID %d set to %f from parent ID %d", task.ID, res, parentID)
				}
			} else {
				ready = false
				log.Printf("CompleteTask: Task ID %d is not ready, waiting for parent task ID %d", task.ID, parentID)
				break
			}
		}

		// Если все зависимости выполнены, добавляем задачу в taskQueue
		if ready {
			// Проверка деления на ноль перед добавлением в очередь
			if task.Operation == "/" && task.Arg2 == 0 {
				log.Printf("CompleteTask: Division by zero detected for task ID %d", task.ID)
				return errors.New("division by zero detected")
			}
			// Логируем изменение флага готовности
			if !task.IsReady {
				task.IsReady = true
				log.Printf("CompleteTask: Task ID %d is now ready (IsReady = %v)", task.ID, task.IsReady)
			}
			app.taskQueue = append(app.taskQueue, task)
			log.Printf("CompleteTask: Task ID %d added to taskQueue", task.ID)
		} else {
			newDependentQueue = append(newDependentQueue, task)
			log.Printf("CompleteTask: Task ID %d is not ready, remaining in dependentQueue", task.ID)
		}
	}

	// Обновляем очередь зависимых задач
	app.dependentQueue = newDependentQueue
	log.Printf("CompleteTask: Updated dependentQueue. Number of tasks remaining: %d", len(app.dependentQueue))

	return nil
}

// GetNextTask выдает агенту следующую задачу
func (app *Application) GetNextTask() (*Task, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	if len(app.taskQueue) == 0 {
		return nil, nil // Очередь пуста
	}

	task := app.taskQueue[0]
	app.taskQueue = app.taskQueue[1:]

	log.Printf("GetNextTask: Issued task ID %d to agent", task.ID)
	return task, nil
}

// getOperationTime возвращает время выполнения операции
func (app *Application) getOperationTime(operation string) int64 {
	var envVar string
	switch operation {
	case "+":
		envVar = "TIME_ADDITION_MS"
	case "-":
		envVar = "TIME_SUBTRACTION_MS"
	case "*":
		envVar = "TIME_MULTIPLICATIONS_MS"
	case "/":
		envVar = "TIME_DIVISIONS_MS"
	default:
		return 1000 // Значение по умолчанию, если операция неизвестна
	}

	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		log.Printf("Warning: %s is not set. Using default value (1000ms)", envVar)
		return 1000
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil || value <= 0 {
		log.Printf("Error: %s is not a valid positive number. Using default value (1000ms)", envVar)
		return 1000
	}

	return value
}

// GetExpressionByID возвращает выражение по ID

func (app *Application) GetExpressionByID(id int) (*Expression, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	expr, exists := app.expressions[id]
	if !exists {
		return nil, fmt.Errorf("expression ID %d not found", id)
	}

	// Проверяем, все ли задачи выполнены
	allCompleted := true
	for _, task := range expr.Tasks {
		if task.Status != "completed" {
			allCompleted = false
			log.Printf("GetExpressionByID: Expression ID %d has incomplete task ID %d with status %s", id, task.ID, task.Status)
			break
		}
	}

	// Если все задачи выполнены, но статус выражения еще не обновлен
	if allCompleted && expr.Status != "completed" {
		// Финальный результат находится в последней задаче
		if len(expr.Tasks) > 0 {
			finalTask := expr.Tasks[len(expr.Tasks)-1]
			expr.Status = "completed"
			expr.Result = finalTask.Result
			log.Printf("GetExpressionByID: Expression ID %d status updated to 'completed' with result %.2f", id, expr.Result)
		} else {
			// Обработка случая, когда в выражении нет задач (хотя такого быть не должно)
			expr.Status = "completed"
			expr.Result = 0
			log.Printf("GetExpressionByID: Warning - Expression ID %d has no tasks but marked as completed", id)
		}
	}

	return expr, nil
}

// GetAllExpressions возвращает все выражения
func (app *Application) GetAllExpressions() map[int]*Expression {
	app.mu.Lock()
	defer app.mu.Unlock()

	return app.expressions
}

// isNumber проверяет, является ли токен числом
func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

// isOperator проверяет, является ли токен оператором
func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

// GetExpressionResult возвращает результат выражения по его ID
func (app *Application) GetExpressionResult(exprID int) (float64, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	expr, exists := app.expressions[exprID]
	if !exists {
		log.Printf("GetExpressionResult: Expression ID %d not found", exprID)
		return 0, errors.New("expression not found")
	}

	// Проверяем, завершены ли все задачи в выражении
	if expr.Status != "completed" {
		log.Printf("GetExpressionResult: Expression ID %d is not completed yet", exprID)
		return 0, errors.New("result not ready")
	}

	log.Printf("GetExpressionResult: Returning result %.2f for expression ID %d", expr.Result, exprID)
	return expr.Result, nil
}
