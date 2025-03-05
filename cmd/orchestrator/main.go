package main

import (
	"fmt"
	"github.com/mpkelevra23/arithmetic-web-service/internal/orchestrator"
	"log"
	"net/http"
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

	// Получаем времена выполнения операций
	opTimes := orchestrator.OperationTimes{
		Addition:       getEnvInt("TIME_ADDITION_MS", 100),
		Subtraction:    getEnvInt("TIME_SUBTRACTION_MS", 100),
		Multiplication: getEnvInt("TIME_MULTIPLICATIONS_MS", 200),
		Division:       getEnvInt("TIME_DIVISIONS_MS", 200),
	}

	// Получаем порт сервера
	port := getEnv("PORT", "8080")

	// Создаем компоненты сервера
	storage := orchestrator.NewStorage()
	parser := orchestrator.NewParser(opTimes)
	server := orchestrator.NewServer(storage, parser)

	// Настраиваем маршруты
	handler := server.SetupRoutes()

	// Запускаем сервер
	log.Printf("Оркестратор запущен на порту %s\n", port)
	log.Printf("Времена операций: сложение=%dms, вычитание=%dms, умножение=%dms, деление=%dms\n",
		opTimes.Addition, opTimes.Subtraction, opTimes.Multiplication, opTimes.Division)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v\n", err)
	}
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
	valueStr := getEnv(key, fmt.Sprintf("%d", defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Ошибка преобразования %s=%s в число, используем значение по умолчанию %d\n",
			key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
