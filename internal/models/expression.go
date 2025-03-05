package models

// Expression представляет арифметическое выражение и его статус
type Expression struct {
	ID       int     `json:"id"`                   // Уникальный идентификатор выражения
	RawExpr  string  `json:"expression,omitempty"` // Исходное строковое выражение
	Status   Status  `json:"status"`               // Текущий статус вычисления
	Result   *string `json:"result,omitempty"`     // Результат вычисления (nil, если не вычислено)
	ErrorMsg string  `json:"error,omitempty"`      // Сообщение об ошибке (если статус ERROR)
}

// ExpressionRequest представляет запрос на вычисление выражения
type ExpressionRequest struct {
	Expression string `json:"expression"` // Строка с выражением
}

// ExpressionResponse представляет ответ на запрос добавления выражения
type ExpressionResponse struct {
	ID int `json:"id"` // Идентификатор добавленного выражения
}

// ExpressionsResponse представляет ответ со списком выражений
type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"` // Список выражений
}

// ExpressionDetailResponse представляет ответ с детальной информацией о выражении
type ExpressionDetailResponse struct {
	Expression Expression `json:"expression"` // Детальная информация о выражении
}
