package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mpkelevra23/arithmetic-web-service/internal/middleware"
	"go.uber.org/zap"
)

// TestLoggingMiddleware проверяет, что LoggingMiddleware корректно обрабатывает запросы.
func TestLoggingMiddleware(t *testing.T) {
	// Инициализация логгера
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer logger.Sync()

	// Тестовый обработчик, который возвращает статус 200
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Имитация задержки обработки
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем тестовый обработчик в middleware
	loggedHandler := middleware.LoggingMiddleware(logger)(testHandler)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос
	loggedHandler.ServeHTTP(w, req)

	// Проверяем статус код ответа
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Result().StatusCode)
	}
}
