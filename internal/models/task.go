package models

// Operation представляет тип операции в задаче
type Operation string

// Константы для типов операций
const (
	OperationAdd      Operation = "ADD"      // Сложение
	OperationSubtract Operation = "SUBTRACT" // Вычитание
	OperationMultiply Operation = "MULTIPLY" // Умножение
	OperationDivide   Operation = "DIVIDE"   // Деление
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

// TaskResponse представляет ответ на запрос задачи для выполнения
type TaskResponse struct {
	Task *Task `json:"task,omitempty"` // Задача (nil, если нет доступных задач)
}

// TaskResultRequest представляет запрос на добавление результата выполненной задачи
type TaskResultRequest struct {
	ID     int     `json:"id"`     // Идентификатор задачи
	Result float64 `json:"result"` // Результат выполнения
}
