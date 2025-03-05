package tests

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mpkelevra23/arithmetic-web-service/errors"
	"github.com/mpkelevra23/arithmetic-web-service/internal/handler"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
)

// TestCalculateHandler проверяет работу HTTP-обработчика CalculateHandler.
func TestCalculateHandler(t *testing.T) {
	// Инициализация логгера
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Logger sync error", zap.Error(err))
		}
	}(logger)

	// Инициализация обработчика
	handlerFunc := handler.CalculateHandler(logger)

	// Определение тестовых случаев
	tests := []struct {
		name           string
		method         string
		payload        interface{}
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Valid Expression",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "1 + 2 * 3"},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"result": "7"},
		},
		{
			name:           "Invalid Characters",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "1 + a"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrInvalidInput},
		},
		{
			name:           "Division by Zero",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "10 / 0"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrDivisionByZero},
		},
		{
			name:           "Missing Expression Field",
			method:         http.MethodPost,
			payload:        map[string]string{"expr": "1 + 2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrMissingField},
		},
		{
			name:           "Empty Expression",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": ""},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrMissingField},
		},
		{
			name:           "Unsupported HTTP Method",
			method:         http.MethodGet, // Используем неподдерживаемый метод
			payload:        nil,            // Нет тела запроса
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   map[string]string{"error": errors.ErrUnsupportedMethod},
		},
		{
			name:           "Malformed JSON",
			method:         http.MethodPost,
			payload:        `{"expression": "1 + 2",`, // Неправильный JSON
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": errors.ErrMalformedJSON},
		},
		{
			name:           "Expression Too Long",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": generateLongExpression(1001)}, // Предполагаем, что 1001 символ превышает лимит
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrTooLongExpression},
		},
		{
			name:           "Trigger Internal Server Error by Header",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "1+1"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Internal Server Error"},
		},
		{
			name:           "Trigger Internal Server Error by Long Expression",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": generateLongExpression(800)},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Expression length triggered server error"},
		},
	}

	for _, tt := range tests {
		tt := tt // Захват переменной для параллельного выполнения
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			// Подготовка тела запроса
			switch payload := tt.payload.(type) {
			case string:
				reqBody = []byte(payload)
			case map[string]string:
				reqBody, err = json.Marshal(payload)
				if err != nil {
					t.Fatalf("Не удалось сериализовать payload: %v", err)
				}
			case nil:
				reqBody = nil
			default:
				t.Fatalf("Неподдерживаемый тип payload: %T", payload)
			}

			// Создание нового HTTP запроса
			req := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.name == "Trigger Internal Server Error by Header" {
				req.Header.Set("X-Trigger-500", "true")
			}

			// Создание ResponseRecorder для записи ответа
			rr := httptest.NewRecorder()

			// Вызов обработчика
			handlerFunc.ServeHTTP(rr, req)

			// Проверка статуса ответа
			if rr.Code != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, rr.Code)
			}

			// Декодирование тела ответа
			var responseBody map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Не удалось декодировать тело ответа: %v", err)
			}

			// Проверка соответствия ожидаемого и фактического тела ответа
			for key, expectedValue := range tt.expectedBody {
				if value, exists := responseBody[key]; !exists || value != expectedValue {
					t.Errorf("Для ключа '%s' ожидалось '%s', получено '%s'", key, expectedValue, value)
				}
			}
		})
	}
}

// generateLongExpression создает строку арифметического выражения заданной длины.
func generateLongExpression(length int) string {
	expression := ""
	for i := 0; i < length; i++ {
		expression += "1"
	}
	return expression
}
