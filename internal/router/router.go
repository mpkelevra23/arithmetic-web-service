package router

import (
	"github.com/mpkelevra23/arithmetic-web-service/internal/handler"
	"github.com/mpkelevra23/arithmetic-web-service/internal/middleware"
	"go.uber.org/zap"
	"net/http"
)

// NewRouter настраивает маршруты и middleware.
func NewRouter(logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	// Обработчики
	calculateHandler := handler.CalculateHandler(logger)
	mux.Handle("/api/v1/calculate", calculateHandler)

	// Применение middleware
	loggedRouter := middleware.LoggingMiddleware(logger)(mux)

	return loggedRouter
}
