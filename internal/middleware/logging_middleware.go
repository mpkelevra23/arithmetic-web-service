package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

// LoggingMiddleware логирует все входящие HTTP-запросы с информацией о методе, URL, статусе и времени обработки.
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Используем обёртку для захвата статус кода ответа
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			logger.Info("Request",
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

// WriteHeader сохраняет статус код ответа и передаёт его далее.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
