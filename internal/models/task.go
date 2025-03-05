package models

// Operation представляет тип операции в задаче
type Operation string

// Константы для типов операций
const (
	OperationAdd      Operation = "ADD"
	OperationSubtract Operation = "SUBTRACT"
	OperationMultiply Operation = "MULTIPLY"
	OperationDivide   Operation = "DIVIDE"
)

// Task представляет задачу на выполнение одной операции
type Task struct {
	ID            int       `json:"id"`                      // Уникальный идентификатор задачи
	ExpressionID  int       `json:"expression_id,omitempty"` // Идентификатор выражения
	Arg1          string    `json:"arg1"`                    // Первый аргумент
	Arg2          string    `json:"arg2"`                    // Второй аргумент
	Operation     Operation `json:"operation"`               // Операция
	OperationTime int       `json:"operation_time"`          // Время выполнения в миллисекундах
	Result        *float64  `json:"result,omitempty"`        // Результат выполнения
	Dependencies  []int     `json:"-"`                       // Зависимости от других задач
	IsReady       bool      `json:"-"`                       // Готовность к выполнению
}

// TaskResponse представляет запрос на добавление задачи
type TaskResponse struct {
	Task *Task `json:"task,omitempty"`
}

// TaskResultRequest представляет запрос на отправку результата задачи
type TaskResultRequest struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}
