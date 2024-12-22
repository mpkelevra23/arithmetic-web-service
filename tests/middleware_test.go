package tests

import (
	"github.com/mpkelevra23/arithmetic-web-service/internal/middleware"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoggingMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Тестовый обработчик, который просто возвращает статус 200
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Имитация задержки
		w.WriteHeader(http.StatusOK)
	})

	// Обертываем тестовый обработчик в middleware
	loggedHandler := middleware.LoggingMiddleware(logger)(testHandler)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос
	loggedHandler.ServeHTTP(w, req)

	// Проверяем статус код
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Result().StatusCode)
	}
}
