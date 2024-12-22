package router

import (
	"net/http"

	"github.com/mpkelevra23/arithmetic-web-service/internal/handler"
	"github.com/mpkelevra23/arithmetic-web-service/internal/middleware"
	"go.uber.org/zap"
)

// NewRouter настраивает маршруты и применяет middleware.
func NewRouter(logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	// Регистрация обработчика для эндпоинта /api/v1/calculate
	mux.Handle("/api/v1/calculate", handler.CalculateHandler(logger))

	// Применение middleware для логирования запросов
	loggedRouter := middleware.LoggingMiddleware(logger)(mux)

	return loggedRouter
}
