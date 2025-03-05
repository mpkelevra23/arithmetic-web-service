package main

import (
	"github.com/mpkelevra23/arithmetic-web-service/internal/agent"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используем переменные окружения системы")
	}

	// Получаем URL оркестратора
	orchestratorURL := getEnv("ORCHESTRATOR_URL", "http://localhost:8080")

	// Получаем вычислительную мощность
	computingPower := getEnvInt("COMPUTING_POWER", 3)

	// Создаем агента
	a := agent.NewAgent(orchestratorURL, computingPower)

	// Запускаем агента
	log.Printf("Агент запущен. URL оркестратора: %s, вычислительная мощность: %d\n",
		orchestratorURL, computingPower)
	a.Start()
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvInt возвращает целочисленное значение переменной окружения или значение по умолчанию
func getEnvInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Ошибка преобразования %s=%s в число, используем значение по умолчанию %d\n",
			key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
