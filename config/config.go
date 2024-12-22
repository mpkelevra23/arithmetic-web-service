// config/config.go
package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config представляет конфигурационные параметры приложения.
type Config struct {
	Port     string
	LogLevel string
}

// LoadConfig загружает конфигурацию из .env файла и переменных окружения.
// Возвращает конфигурацию, булевый флаг наличия .env файла и ошибку, если она произошла.
func LoadConfig() (*Config, bool, error) {
	// Загрузка переменных из .env файла
	err := godotenv.Load()
	envLoaded := true
	if err != nil {
		// Если файл .env не найден, считаем, что переменные окружения уже установлены
		envLoaded = false
	}

	config := &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return config, envLoaded, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию.
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
