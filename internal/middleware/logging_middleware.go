// internal/middleware/logging_middleware.go
package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

// LoggingMiddleware логирует все входящие HTTP-запросы.
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Используем ResponseWriter, который захватывает статус код
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			logger.Info("Входящий запрос",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", lrw.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}

// loggingResponseWriter оборачивает http.ResponseWriter для захвата статус кода ответа.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader захватывает статус код.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
