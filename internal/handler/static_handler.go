package handler

import (
	"net/http"

	"github.com/mpkelevra23/arithmetic-web-service/web"
	"go.uber.org/zap"
)

// StaticHandler возвращает обработчик для статических файлов фронтенда
func StaticHandler(logger *zap.Logger) http.Handler {
	// Получаем файловую систему для статических файлов
	fileSystem, err := web.GetFileSystem()
	if err != nil {
		logger.Error("Failed to get web filesystem", zap.Error(err))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		})
	}

	// Создаем файловый сервер
	fileServer := http.FileServer(fileSystem)

	// Обрабатываем запросы
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отладочное логирование
		logger.Debug("Static request", zap.String("path", r.URL.Path))

		// Если запрос на корневой путь, перенаправляем на index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}

		// Обслуживаем запрошенный файл
		fileServer.ServeHTTP(w, r)
	})
}
