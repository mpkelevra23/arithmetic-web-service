package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/mpkelevra23/arithmetic-web-service/errors"
	"github.com/mpkelevra23/arithmetic-web-service/internal/calculator"
	"go.uber.org/zap"
)

// CalculateRequest представляет структуру входящего запроса.
type CalculateRequest struct {
	Expression string `json:"expression"`
}

// CalculateResponse представляет структуру успешного ответа.
type CalculateResponse struct {
	Result string `json:"result"`
}

// expressionRegex используется для валидации допустимых символов в выражении.
var expressionRegex = regexp.MustCompile(`^[0-9+\-*/().\s]+$`)

// CalculateHandler обрабатывает POST-запросы к эндпоинту /api/v1/calculate.
func CalculateHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CalculateRequest

		// Пример искусственного вызова ошибки 500
		if r.Header.Get("X-Trigger-500") == "true" {
			logger.Error("Artificial internal server error triggered")
			errors.WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// Проверка метода запроса
		if r.Method != http.MethodPost {
			logger.Warn("Unsupported HTTP method", zap.String("method", r.Method))
			errors.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.ErrUnsupportedMethod)
			return
		}

		// Декодирование JSON тела запроса
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Request body decoding error", zap.Error(err))
			errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrMalformedJSON)
			return
		}

		// Очистка и проверка поля expression
		expression := strings.TrimSpace(req.Expression)
		if expression == "" {
			logger.Error("Empty expression field")
			errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrMissingField)
			return
		}

		// Валидация выражения на допустимые символы
		if !expressionRegex.MatchString(expression) {
			logger.Error("Invalid characters in expression", zap.String("expression", expression))
			errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrInvalidInput)
			return
		}

		// Пример искусственного вызова ошибки 500 на основе длины выражения
		if len(expression) > 500 && len(expression) <= 1000 {
			logger.Error("Artificial internal server error for long expression")
			errors.WriteErrorResponse(w, http.StatusInternalServerError, "Expression length triggered server error")
			return
		}

		// Проверка длины выражения
		if len(expression) > 1000 {
			logger.Error("Expression is too long", zap.Int("length", len(expression)))
			errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrTooLongExpression)
			return
		}

		// Вычисление результата выражения
		result, err := calculator.Calc(expression)
		if err != nil {
			logger.Error("Calculation error", zap.Error(err))
			handleCalculationError(w, err)
			return
		}

		// Формирование успешного ответа
		resp := CalculateResponse{
			Result: formatResult(result),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("Response encoding error", zap.Error(err))
		}
	}
}

// handleCalculationError обрабатывает ошибки, возникшие при вычислении выражения.
func handleCalculationError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "division by zero":
		errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrDivisionByZero)
	default:
		errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrInvalidExpression)
	}
}

// formatResult форматирует результат вычисления, убирая лишние нули.
func formatResult(result float64) string {
	if result == float64(int64(result)) {
		return strconv.FormatInt(int64(result), 10)
	}
	return strconv.FormatFloat(result, 'f', -1, 64)
}
