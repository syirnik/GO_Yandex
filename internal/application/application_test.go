package application

import (
	"sync"
	"testing"
	"time"
)

// Тест инициализации приложения
func TestNew(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("New() returned nil")
	}

	if app.nextTaskID != 1 {
		t.Errorf("Expected nextTaskID to be 1, got %d", app.nextTaskID)
	}

	if app.nextExpressionID != 1 {
		t.Errorf("Expected nextExpressionID to be 1, got %d", app.nextExpressionID)
	}

	if len(app.expressions) != 0 {
		t.Errorf("Expected empty expressions map, got %d items", len(app.expressions))
	}

	if len(app.taskQueue) != 0 {
		t.Errorf("Expected empty taskQueue, got %d items", len(app.taskQueue))
	}
}

// Тест разбора простого выражения
func TestParseSimpleExpression(t *testing.T) {
	app := New()
	exprID, err := app.ParseExpression("2+3")

	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	if exprID != 1 {
		t.Errorf("Expected expression ID 1, got %d", exprID)
	}

	expr, err := app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Value != "2+3" {
		t.Errorf("Expected expression value '2+3', got '%s'", expr.Value)
	}

	if expr.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", expr.Status)
	}

	if len(expr.Tasks) == 0 {
		t.Error("Expected tasks to be created, but none found")
	}
}

// Тест разбора сложного выражения
func TestParseComplexExpression(t *testing.T) {
	app := New()
	exprID, err := app.ParseExpression("2+3*4")

	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	expr, err := app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Value != "2+3*4" {
		t.Errorf("Expected expression value '2+3*4', got '%s'", expr.Value)
	}

	if len(expr.Tasks) == 0 {
		t.Error("Expected tasks to be created, but none found")
	}
}

// Тест обработки ошибок при разборе выражений
func TestParseExpressionErrors(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
	}{
		{"Division By Zero", "5/0"},
		{"Invalid Expression", "2++3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := New()
			_, err := app.ParseExpression(tc.expression)

			if err == nil {
				t.Errorf("Expected error for expression '%s', got nil", tc.expression)
			}
		})
	}
}

// Тест получения и выполнения задач
func TestTaskProcessing(t *testing.T) {
	app := New()
	exprID, err := app.ParseExpression("2+3")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	// Получаем задачу
	task, err := app.GetNextTask()
	if err != nil {
		t.Fatalf("GetNextTask returned error: %v", err)
	}

	if task == nil {
		t.Fatal("Expected task, got nil")
	}

	if task.Operation != "+" {
		t.Errorf("Expected operation '+', got '%s'", task.Operation)
	}

	// Выполняем задачу
	err = app.CompleteTask(task.ID, 5.0)
	if err != nil {
		t.Fatalf("CompleteTask returned error: %v", err)
	}

	// Проверяем, что результат сохранен
	expr, err := app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", expr.Status)
	}

	if expr.Result != 5.0 {
		t.Errorf("Expected result 5.0, got %f", expr.Result)
	}
}

// Тест обработки зависимостей между задачами
func TestTaskDependencies(t *testing.T) {
	app := New()
	_, err := app.ParseExpression("2+3*4")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	// Выполняем все задачи
	processAllTasks(t, app)

	// Проверяем конечный результат
	expr, err := app.GetExpressionByID(1)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", expr.Status)
	}

	// 2+3*4 = 2+12 = 14
	if expr.Result != 14.0 {
		t.Errorf("Expected result 14.0, got %f", expr.Result)
	}
}

// Тест получения информации о выражениях
func TestGetExpressions(t *testing.T) {
	app := New()
	id1, _ := app.ParseExpression("2+3")
	id2, _ := app.ParseExpression("4*5")

	// Получение по ID
	expr1, err := app.GetExpressionByID(id1)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr1.Value != "2+3" {
		t.Errorf("Expected expression value '2+3', got '%s'", expr1.Value)
	}

	// Получение всех выражений
	expressions := app.GetAllExpressions()
	if len(expressions) != 2 {
		t.Errorf("Expected 2 expressions, got %d", len(expressions))
	}

	if _, exists := expressions[id1]; !exists {
		t.Errorf("Expression with ID %d not found", id1)
	}

	if _, exists := expressions[id2]; !exists {
		t.Errorf("Expression with ID %d not found", id2)
	}

	// Получение несуществующего выражения
	_, err = app.GetExpressionByID(999)
	if err == nil {
		t.Error("Expected error for non-existent expression ID, got nil")
	}
}

// Тест вспомогательных функций
func TestHelperFunctions(t *testing.T) {
	// Тест isNumber
	numbers := []string{"123", "0.123", "-123.456"}
	for _, num := range numbers {
		if !isNumber(num) {
			t.Errorf("Expected '%s' to be recognized as a number", num)
		}
	}

	notNumbers := []string{"abc", "", "+"}
	for _, notNum := range notNumbers {
		if isNumber(notNum) {
			t.Errorf("Expected '%s' to not be recognized as a number", notNum)
		}
	}

	// Тест isOperator
	operators := []string{"+", "-", "*", "/"}
	for _, op := range operators {
		if !isOperator(op) {
			t.Errorf("Expected '%s' to be recognized as an operator", op)
		}
	}

	notOperators := []string{"123", "abc", ""}
	for _, notOp := range notOperators {
		if isOperator(notOp) {
			t.Errorf("Expected '%s' to not be recognized as an operator", notOp)
		}
	}
}

// Тест многопоточности
func TestConcurrency(t *testing.T) {
	app := New()

	// Создаем несколько выражений
	numExpressions := 3
	exprIDs := make([]int, numExpressions)

	for i := 0; i < numExpressions; i++ {
		id, err := app.ParseExpression("2+3")
		if err != nil {
			t.Fatalf("ParseExpression returned error: %v", err)
		}
		exprIDs[i] = id
	}

	// Запускаем несколько горутин для обработки задач
	var wg sync.WaitGroup
	numWorkers := 2
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()

			for {
				task, err := app.GetNextTask()
				if err != nil {
					t.Errorf("GetNextTask returned error: %v", err)
					return
				}

				if task == nil {
					return // Нет больше задач
				}

				// Симулируем время выполнения
				time.Sleep(time.Millisecond * 10)

				// Выполняем задачу
				result := task.Arg1 + task.Arg2
				err = app.CompleteTask(task.ID, result)
				if err != nil {
					t.Errorf("CompleteTask returned error: %v", err)
					return
				}
			}
		}()
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Проверяем, что все выражения выполнены
	for _, id := range exprIDs {
		expr, err := app.GetExpressionByID(id)
		if err != nil {
			t.Fatalf("GetExpressionByID returned error: %v", err)
		}

		if expr.Status != "completed" {
			t.Errorf("Expression ID %d: expected status 'completed', got '%s'", id, expr.Status)
		}

		if expr.Result != 5.0 {
			t.Errorf("Expression ID %d: expected result 5.0, got %f", id, expr.Result)
		}
	}
}

// Тест различных типов выражений
func TestVariousExpressions(t *testing.T) {
	testCases := []struct {
		expression string
		expected   float64
	}{
		{"2+3", 5.0},
		{"4*5", 20.0},
		{"10-5", 5.0},
		{"8/2", 4.0},
		{"2+3*4", 14.0},
		{"(2+3)*4", 20.0},
	}

	for _, tc := range testCases {
		t.Run(tc.expression, func(t *testing.T) {
			app := New()
			exprID, err := app.ParseExpression(tc.expression)

			if err != nil {
				t.Fatalf("ParseExpression('%s') returned error: %v", tc.expression, err)
			}

			// Обрабатываем все задачи
			processAllTasks(t, app)

			// Проверяем результат
			expr, err := app.GetExpressionByID(exprID)
			if err != nil {
				t.Fatalf("GetExpressionByID returned error: %v", err)
			}

			if expr.Status != "completed" {
				t.Errorf("Expected status 'completed', got '%s'", expr.Status)
			}

			if expr.Result != tc.expected {
				t.Errorf("Expression '%s': expected result %f, got %f", tc.expression, tc.expected, expr.Result)
			}
		})
	}
}

// Вспомогательная функция для обработки всех задач
func processAllTasks(t *testing.T, app *Application) {
	for {
		task, err := app.GetNextTask()
		if err != nil {
			t.Fatalf("GetNextTask returned error: %v", err)
		}

		if task == nil {
			break // Нет больше задач
		}

		var result float64
		switch task.Operation {
		case "+":
			result = task.Arg1 + task.Arg2
		case "-":
			result = task.Arg1 - task.Arg2
		case "*":
			result = task.Arg1 * task.Arg2
		case "/":
			if task.Arg2 == 0 {
				t.Fatalf("Division by zero detected in task ID %d", task.ID)
			}
			result = task.Arg1 / task.Arg2
		}

		err = app.CompleteTask(task.ID, result)
		if err != nil {
			t.Fatalf("CompleteTask returned error: %v", err)
		}
	}
}

// Тест для проверки полного цикла обработки выражения
func TestFullExpressionProcessing(t *testing.T) {
	app := New()
	exprID, err := app.ParseExpression("(2+3)*(4-1)")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	// Обрабатываем все задачи
	processAllTasks(t, app)

	// Проверяем результат
	expr, err := app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", expr.Status)
	}

	// (2+3)*(4-1) = 5*3 = 15
	if expr.Result != 15.0 {
		t.Errorf("Expected result 15.0, got %f", expr.Result)
	}
}

// Тест обработки выражения с несколькими операциями
func TestMultipleOperationsExpression(t *testing.T) {
	app := New()
	exprID, err := app.ParseExpression("2+3*4-5")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	// Обрабатываем все задачи
	processAllTasks(t, app)

	// Проверяем результат
	expr, err := app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", expr.Status)
	}

	// 2+3*4-5 = 2+12-5 = 9
	if expr.Result != 9.0 {
		t.Errorf("Expected result 9.0, got %f", expr.Result)
	}
}

// Тест для проверки обработки вложенных скобок
func TestNestedParentheses(t *testing.T) {
	app := New()
	exprID, err := app.ParseExpression("(2+(3*4))-5")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	// Обрабатываем все задачи
	processAllTasks(t, app)

	// Проверяем результат
	expr, err := app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", expr.Status)
	}

	// (2+(3*4))-5 = (2+12)-5 = 14-5 = 9
	if expr.Result != 9.0 {
		t.Errorf("Expected result 9.0, got %f", expr.Result)
	}
}

// Тест для проверки обновления состояния выражения
func TestExpressionStatusUpdate(t *testing.T) {
	app := New()
	exprID, err := app.ParseExpression("2+3")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	// Изначально статус должен быть pending
	expr, err := app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Status != "pending" {
		t.Errorf("Initial status: expected 'pending', got '%s'", expr.Status)
	}

	// Выполняем задачу
	task, _ := app.GetNextTask()
	err = app.CompleteTask(task.ID, 5.0)
	if err != nil {
		t.Fatalf("CompleteTask returned error: %v", err)
	}

	// После выполнения статус должен быть completed
	expr, err = app.GetExpressionByID(exprID)
	if err != nil {
		t.Fatalf("GetExpressionByID returned error: %v", err)
	}

	if expr.Status != "completed" {
		t.Errorf("Final status: expected 'completed', got '%s'", expr.Status)
	}

	if expr.Result != 5.0 {
		t.Errorf("Expected result 5.0, got %f", expr.Result)
	}
}
