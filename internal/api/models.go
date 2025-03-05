package api

// RequestAddExpression представляет тело запроса для добавления выражения
type RequestAddExpression struct {
	Expression string `json:"expression"`
}

// ResponseAddExpression представляет тело ответа для добавления выражения
type ResponseAddExpression struct {
	ID int `json:"id"`
}

// TaskResponse представляет данные задачи
type TaskResponse struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
}

// ResponseGetTask представляет тело ответа для получения задачи
type ResponseGetTask struct {
	Task TaskResponse `json:"task"`
}

// RequestPostTask представляет тело запроса для завершения задачи
type RequestPostTask struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}

// GetExpressionResponse представляет тело ответа для получения выражения по ID
type GetExpressionResponse struct {
	Expression ExpressionResponse `json:"expression"`
}

// ExpressionResponse представляет данные одного выражения
type ExpressionResponse struct {
	ID     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// ResponseGetExpressions представляет тело ответа для получения списка выражений
type ResponseGetExpressions struct {
	Expressions []ExpressionResponse `json:"expressions"`
}
