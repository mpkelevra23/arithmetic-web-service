package main

import (
	"net/http"
	"os"

	"github.com/mpkelevra23/arithmetic-web-service/config"
	"github.com/mpkelevra23/arithmetic-web-service/internal/router"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Загрузка конфигурации
	cfg, envLoaded, err := config.LoadConfig()
	if err != nil {
		zapLogger, _ := zap.NewProduction()
		defer zapLogger.Sync()
		zapLogger.Fatal("Configuration error", zap.Error(err))
	}

	// Инициализация логгера с указанным уровнем логирования
	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		zapLogger, _ := zap.NewProduction()
		defer zapLogger.Sync()
		zapLogger.Fatal("Logger initialization error", zap.Error(err))
	}
	defer logger.Sync()

	// Информация о загрузке .env файла
	if envLoaded {
		logger.Info(".env file successfully loaded")
	} else {
		logger.Info(".env file not found")
	}

	// Создание маршрутизатора
	r := router.NewRouter(logger)

	// Запуск HTTP-сервера
	logger.Info("Server started", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

// initLogger инициализирует логгер Zap с заданным уровнем логирования.
func initLogger(level string) (*zap.Logger, error) {

	// Убедимся, что директория для логов существует
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	
	// Настройка конфигурации логгера
	zapConfig := zap.NewProductionConfig()

	// Установка уровня логирования
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

	// Настройка путей вывода логов
	zapConfig.OutputPaths = []string{
		"stdout",
		"logs/app.log",
	}
	zapConfig.ErrorOutputPaths = []string{
		"stderr",
		"logs/error.log",
	}

	// Конфигурация формата времени в логах
	zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Создание логгера
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
