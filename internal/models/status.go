package models

// Status определяет текущий статус выражения
type Status string

// Константы для возможных статусов выражения
const (
	StatusPending    Status = "PENDING"    // Ожидает выполнения
	StatusProcessing Status = "PROCESSING" // В процессе выполнения
	StatusCompleted  Status = "COMPLETED"  // Вычисление завершено
	StatusError      Status = "ERROR"      // Ошибка при вычислении
)
