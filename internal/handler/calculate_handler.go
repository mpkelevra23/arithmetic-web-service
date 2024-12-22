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

// Предкомпилированное регулярное выражение для валидации.
var expressionRegex = regexp.MustCompile(`^[0-9+\-*/().\s]+$`)

// CalculateHandler обрабатывает запросы к эндпоинту /api/v1/calculate.
func CalculateHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CalculateRequest

		// Ограничение метода запроса
		if r.Method != http.MethodPost {
			logger.Warn("Unsupported HTTP method", zap.String("method", r.Method))
			errors.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.ErrUnsupportedMethod)
			return
		}

		// Декодирование JSON-запроса
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Ошибка декодирования запроса", zap.Error(err))
			errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrMalformedJSON)
			return
		}

		// Проверка наличия поля expression
		expression := strings.TrimSpace(req.Expression)
		if expression == "" {
			logger.Error("Отсутствует поле expression")
			errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrMissingField)
			return
		}

		// Валидация выражения на допустимые символы
		if !expressionRegex.MatchString(expression) {
			logger.Error("Недопустимые символы во входном выражении", zap.String("expression", expression))
			errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrInvalidInput)
			return
		}

		// Дополнительная проверка на длину выражения
		if len(expression) > 1000 {
			logger.Error("Выражение слишком длинное", zap.Int("length", len(expression)))
			errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrTooLongExpression)
			return
		}

		// Вычисление результата
		result, err := calculator.Calc(expression)
		if err != nil {
			logger.Error("Ошибка вычисления выражения", zap.Error(err))
			switch err.Error() {
			case "деление на ноль":
				errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrDivisionByZero)
			default:
				errors.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.ErrInvalidExpression)
			}
			return
		}

		// Формирование успешного ответа
		resp := CalculateResponse{
			Result: formatResult(result),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("Ошибка кодирования ответа", zap.Error(err))
		}
	}
}

// formatResult форматирует результат вычисления, убирая лишние нули.
func formatResult(result float64) string {
	// Если число целое, представляем без десятичной части
	if result == float64(int64(result)) {
		return strconv.FormatInt(int64(result), 10)
	}
	return strconv.FormatFloat(result, 'f', -1, 64)
}
