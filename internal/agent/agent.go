package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mpkelevra23/arithmetic-web-service/internal/models"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Agent представляет агента, выполняющего задачи
type Agent struct {
	orchestratorURL string
	computingPower  int
	client          *http.Client
	wg              sync.WaitGroup
}

// NewAgent создает нового агента
func NewAgent(orchestratorURL string, computingPower int) *Agent {
	return &Agent{
		orchestratorURL: orchestratorURL,
		computingPower:  computingPower,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Start запускает агента
func (a *Agent) Start() {
	log.Printf("Запуск агента с %d воркерами\n", a.computingPower)

	// Запускаем воркеры
	for i := 0; i < a.computingPower; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}

	// Ожидаем завершения всех воркеров
	a.wg.Wait()
}

// worker представляет горутину, выполняющую задачи
func (a *Agent) worker(id int) {
	defer a.wg.Done()

	log.Printf("Воркер %d запущен\n", id)

	for {
		// Запрашиваем задачу
		task, err := a.getTask()
		if err != nil {
			log.Printf("Воркер %d: ошибка получения задачи: %v\n", id, err)
			time.Sleep(1 * time.Second) // Пауза перед следующей попыткой
			continue
		}

		log.Printf("Воркер %d: получена задача %d (%s %s %s)\n", id, task.ID, task.Arg1, task.Operation, task.Arg2)

		// Выполняем задачу
		result, err := a.executeTask(task)
		if err != nil {
			log.Printf("Воркер %d: ошибка выполнения задачи %d: %v\n", id, task.ID, err)
			continue
		}

		log.Printf("Воркер %d: задача %d выполнена, результат: %f\n", id, task.ID, result)

		// Отправляем результат
		if err := a.sendResult(task.ID, result); err != nil {
			log.Printf("Воркер %d: ошибка отправки результата задачи %d: %v\n", id, task.ID, err)
		}
	}
}

// getTask запрашивает задачу у оркестратора
func (a *Agent) getTask() (*models.Task, error) {
	resp, err := a.client.Get(fmt.Sprintf("%s/internal/task", a.orchestratorURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("нет доступных задач")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный код ответа: %d", resp.StatusCode)
	}

	var taskResp models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}

	return taskResp.Task, nil
}

// executeTask выполняет задачу и возвращает результат
func (a *Agent) executeTask(task *models.Task) (float64, error) {
	// Парсим аргументы
	arg1, err := strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return 0, fmt.Errorf("некорректный аргумент 1: %s", task.Arg1)
	}

	arg2, err := strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return 0, fmt.Errorf("некорректный аргумент 2: %s", task.Arg2)
	}

	// Искусственная задержка
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	// Выполняем операцию
	var result float64
	switch task.Operation {
	case models.OperationAdd:
		result = arg1 + arg2
	case models.OperationSubtract:
		result = arg1 - arg2
	case models.OperationMultiply:
		result = arg1 * arg2
	case models.OperationDivide:
		if arg2 == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		result = arg1 / arg2
	default:
		return 0, fmt.Errorf("неизвестная операция: %s", task.Operation)
	}

	return result, nil
}

// sendResult отправляет результат задачи оркестратору
func (a *Agent) sendResult(taskID int, result float64) error {
	reqBody := models.TaskResultRequest{
		ID:     taskID,
		Result: result,
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := a.client.Post(
		fmt.Sprintf("%s/internal/task", a.orchestratorURL),
		"application/json",
		bytes.NewBuffer(reqData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный код ответа: %d", resp.StatusCode)
	}

	return nil
}
