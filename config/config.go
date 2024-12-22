package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config содержит конфигурационные параметры приложения.
type Config struct {
	Port     string
	LogLevel string
}

// LoadConfig загружает конфигурацию из файла .env и переменных окружения.
// Возвращает конфигурацию, флаг наличия .env файла и ошибку при необходимости.
func LoadConfig() (*Config, bool, error) {
	// Попытка загрузить переменные из .env файла
	err := godotenv.Load()
	envLoaded := err == nil

	config := &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return config, envLoaded, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию, если переменная не установлена.
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
