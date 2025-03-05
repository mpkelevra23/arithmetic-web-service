package agent

import (
	"github.com/mpkelevra23/arithmetic-web-service/internal/models"
	"testing"
	"time"
)

// TestExecuteTask проверяет выполнение задач
func TestExecuteTask(t *testing.T) {
	// Создаем агента
	a := NewAgent("http://localhost:8080", 1)

	// Тесты для разных операций
	tests := []struct {
		name      string
		task      models.Task
		want      float64
		wantError bool
	}{
		{
			name: "Сложение",
			task: models.Task{
				ID:            1,
				Arg1:          "2",
				Arg2:          "3",
				Operation:     models.OperationAdd,
				OperationTime: 1, // Минимальное время для быстрого теста
			},
			want:      5,
			wantError: false,
		},
		{
			name: "Вычитание",
			task: models.Task{
				ID:            2,
				Arg1:          "5",
				Arg2:          "3",
				Operation:     models.OperationSubtract,
				OperationTime: 1,
			},
			want:      2,
			wantError: false,
		},
		{
			name: "Умножение",
			task: models.Task{
				ID:            3,
				Arg1:          "2",
				Arg2:          "3",
				Operation:     models.OperationMultiply,
				OperationTime: 1,
			},
			want:      6,
			wantError: false,
		},
		{
			name: "Деление",
			task: models.Task{
				ID:            4,
				Arg1:          "6",
				Arg2:          "3",
				Operation:     models.OperationDivide,
				OperationTime: 1,
			},
			want:      2,
			wantError: false,
		},
		{
			name: "Деление на ноль",
			task: models.Task{
				ID:            5,
				Arg1:          "6",
				Arg2:          "0",
				Operation:     models.OperationDivide,
				OperationTime: 1,
			},
			want:      0, // Значение не важно, т.к. ожидается ошибка
			wantError: true,
		},
		{
			name: "Некорректный аргумент",
			task: models.Task{
				ID:            6,
				Arg1:          "abc",
				Arg2:          "3",
				Operation:     models.OperationAdd,
				OperationTime: 1,
			},
			want:      0, // Значение не важно, т.к. ожидается ошибка
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Замеряем время начала
			start := time.Now()

			// Выполняем задачу
			got, err := a.executeTask(&tt.task)

			// Проверяем время выполнения
			elapsed := time.Since(start)
			if elapsed < time.Duration(tt.task.OperationTime)*time.Millisecond {
				t.Errorf("executeTask() took %v, want at least %v", elapsed, time.Duration(tt.task.OperationTime)*time.Millisecond)
			}

			// Проверяем наличие ошибки
			if (err != nil) != tt.wantError {
				t.Errorf("executeTask() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// Если ошибка ожидается, не проверяем результат
			if tt.wantError {
				return
			}

			// Проверяем результат
			if got != tt.want {
				t.Errorf("executeTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
