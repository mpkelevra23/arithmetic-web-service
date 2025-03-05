package orchestrator

import (
	"fmt"
	"github.com/mpkelevra23/arithmetic-web-service/internal/models"
	"sync"
)

// Storage представляет in-memory хранилище для выражений и задач
type Storage struct {
	expressions      map[int]models.Expression // Хранилище выражений
	tasks            map[int]models.Task       // Хранилище задач
	exprCounter      int                       // Счетчик для ID выражений
	taskCounter      int                       // Счетчик для ID задач
	mutex            sync.RWMutex              // Мьютекс для защиты данных
	exprTasksMapping map[int][]int             // Связь выражений с задачами
	resultCache      map[string]float64        // Кеш результатов задач
}

// NewStorage создает новое хранилище
func NewStorage() *Storage {
	return &Storage{
		expressions:      make(map[int]models.Expression),
		tasks:            make(map[int]models.Task),
		exprCounter:      0,
		taskCounter:      0,
		exprTasksMapping: make(map[int][]int),
		resultCache:      make(map[string]float64),
	}
}

// AddExpression добавляет новое выражение в хранилище
func (s *Storage) AddExpression(expr string) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.exprCounter++
	id := s.exprCounter

	s.expressions[id] = models.Expression{
		ID:      id,
		RawExpr: expr,
		Status:  models.StatusPending,
	}

	return id, nil
}

// GetExpression возвращает выражение по ID
func (s *Storage) GetExpression(id int) (models.Expression, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	expr, exists := s.expressions[id]
	if !exists {
		return models.Expression{}, fmt.Errorf("выражение с ID %d не найдено", id)
	}

	return expr, nil
}

// GetAllExpressions возвращает все выражения
func (s *Storage) GetAllExpressions() []models.Expression {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]models.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		result = append(result, expr)
	}

	return result
}

// AddTasks добавляет задачи для выражения
func (s *Storage) AddTasks(exprID int, tasks []models.Task) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.expressions[exprID]
	if !exists {
		return fmt.Errorf("выражение с ID %d не найдено", exprID)
	}

	taskIDs := make([]int, 0, len(tasks))

	for i := range tasks {
		s.taskCounter++
		tasks[i].ID = s.taskCounter
		tasks[i].ExpressionID = exprID

		// Задача готова, если нет зависимостей
		tasks[i].IsReady = len(tasks[i].Dependencies) == 0

		s.tasks[tasks[i].ID] = tasks[i]
		taskIDs = append(taskIDs, tasks[i].ID)
	}

	s.exprTasksMapping[exprID] = taskIDs

	// Обновляем статус выражения
	expr := s.expressions[exprID]
	expr.Status = models.StatusProcessing
	s.expressions[exprID] = expr

	return nil
}

// GetReadyTask возвращает задачу, готовую к выполнению
func (s *Storage) GetReadyTask() (*models.Task, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, task := range s.tasks {
		if task.IsReady && task.Result == nil {
			// Помечаем задачу как "в процессе"
			task.IsReady = false
			s.tasks[id] = task

			// Копируем задачу для возврата
			taskToReturn := task

			// Заменяем ссылки на результаты на их значения
			if len(taskToReturn.Arg1) > 4 && taskToReturn.Arg1[:4] == "res:" {
				resultName := taskToReturn.Arg1
				if val, exists := s.resultCache[resultName]; exists {
					taskToReturn.Arg1 = fmt.Sprintf("%f", val)
				}
			}

			if len(taskToReturn.Arg2) > 4 && taskToReturn.Arg2[:4] == "res:" {
				resultName := taskToReturn.Arg2
				if val, exists := s.resultCache[resultName]; exists {
					taskToReturn.Arg2 = fmt.Sprintf("%f", val)
				}
			}

			return &taskToReturn, nil
		}
	}

	return nil, fmt.Errorf("нет готовых задач")
}

// UpdateTaskResult обновляет результат выполненной задачи
func (s *Storage) UpdateTaskResult(id int, result float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	// Обновляем результат задачи
	resultValue := result
	task.Result = &resultValue
	s.tasks[id] = task

	// Добавляем результат в кеш
	resultName := fmt.Sprintf("res:%d", id)
	s.resultCache[resultName] = result

	// Обновляем зависимости других задач
	s.updateDependencies(id)

	// Проверяем завершение выражения
	s.checkExpressionCompletion(task.ExpressionID)

	return nil
}

// updateDependencies обновляет зависимости задач
func (s *Storage) updateDependencies(completedTaskID int) {
	for id, task := range s.tasks {
		if task.Result == nil { // Только невыполненные задачи
			for i, depID := range task.Dependencies {
				if depID == completedTaskID {
					// Удаляем выполненную зависимость
					task.Dependencies = append(task.Dependencies[:i], task.Dependencies[i+1:]...)
					break
				}
			}

			// Если зависимостей нет, задача готова
			if len(task.Dependencies) == 0 {
				task.IsReady = true
			}

			s.tasks[id] = task
		}
	}
}

// checkExpressionCompletion проверяет завершение выражения
func (s *Storage) checkExpressionCompletion(exprID int) {
	taskIDs, exists := s.exprTasksMapping[exprID]
	if !exists {
		return
	}

	allCompleted := true
	var finalResult *float64

	for _, taskID := range taskIDs {
		task, exists := s.tasks[taskID]
		if !exists || task.Result == nil {
			allCompleted = false
			break
		}

		// Последний результат будет финальным
		finalResult = task.Result
	}

	if allCompleted && finalResult != nil {
		expr, exists := s.expressions[exprID]
		if exists {
			expr.Status = models.StatusCompleted
			resultStr := fmt.Sprintf("%g", *finalResult)
			expr.Result = &resultStr
			s.expressions[exprID] = expr
		}
	}
}
