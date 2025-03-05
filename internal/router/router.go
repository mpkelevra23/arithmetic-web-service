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

	// Регистрация обработчика для статических файлов
	mux.Handle("/static/", http.StripPrefix("/static/", handler.StaticHandler(logger)))

	// Регистрация обработчика для корневого пути
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/static/index.html", http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	}))

	// Применение middleware для логирования запросов
	loggedRouter := middleware.LoggingMiddleware(logger)(mux)

	return loggedRouter
}
