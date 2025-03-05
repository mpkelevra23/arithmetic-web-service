package models

// Expression определяет текущий статус выражения
type Expression struct {
	ID       int     `json:"id"`                   // Уникальный идентификатор выражения
	RawExpr  string  `json:"expression,omitempty"` // Исходное строковое выражение
	Status   Status  `json:"status"`               // Текущий статус вычисления
	Result   *string `json:"result,omitempty"`     // Результат вычисления (nil, если не вычислено)
	ErrorMsg string  `json:"error,omitempty"`      // Сообщение об ошибке (если статус ERROR)
}

// ExpressionRequest представляет запрос на добавление выражения
type ExpressionRequest struct {
	Expression string `json:"expression"`
}

// ExpressionResponse представляет ответ с ID добавленного выражения
type ExpressionResponse struct {
	ID int `json:"id"`
}

// ExpressionsResponse представляет ответ со списком выражений
type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

// ExpressionDetailResponse представляет ответ с деталями выражения
type ExpressionDetailResponse struct {
	Expression Expression `json:"expression"`
}
