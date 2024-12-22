// cmd/server/main.go
package main

import (
	"github.com/mpkelevra23/arithmetic-web-service/config"
	"github.com/mpkelevra23/arithmetic-web-service/internal/router"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

func main() {
	// Загрузка конфигурации
	cfg, envLoaded, err := config.LoadConfig()
	if err != nil {
		// Если произошла ошибка при загрузке конфигурации, логируем и завершаем работу
		// Используем Zap вместо log.Fatalf
		zapLogger, _ := zap.NewProduction()
		defer zapLogger.Sync()
		zapLogger.Fatal("Ошибка загрузки конфигурации", zap.Error(err))
	}

	// Настройка логирования Zap
	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		zapLogger, _ := zap.NewProduction()
		defer zapLogger.Sync()
		zapLogger.Fatal("Ошибка инициализации логгера", zap.Error(err))
	}
	defer logger.Sync()

	// Логируем информацию о загрузке .env файла
	if envLoaded {
		logger.Info("Файл .env успешно загружен")
	} else {
		logger.Info("Файл .env не найден, используется переменные окружения")
	}

	// Инициализация роутера
	r := router.NewRouter(logger)

	// Запуск сервера
	logger.Info("Запуск сервера", zap.String("порт", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatal("Сервер завершил работу с ошибкой", zap.Error(err))
	}
}

// initLogger инициализирует логгер Zap с заданным уровнем логирования.
func initLogger(level string) (*zap.Logger, error) {
	var zapConfig zap.Config

	// Используем конфигурацию по умолчанию
	zapConfig = zap.NewProductionConfig()

	// Настройка уровня логирования
	switch level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Настройка вывода в консоль и файл
	zapConfig.OutputPaths = []string{
		"stdout",
		"logs/app.log",
	}
	zapConfig.ErrorOutputPaths = []string{
		"stderr",
		"logs/error.log",
	}

	// Изменяем формат времени
    zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
    zapConfig.EncoderConfig.TimeKey = "timestamp" // Измените название ключа времени, если нужно
    zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Создание логгера
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
