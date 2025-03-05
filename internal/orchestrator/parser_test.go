package orchestrator

import (
	"github.com/mpkelevra23/arithmetic-web-service/internal/models"
	"testing"
)

func TestParser_ParseExpression(t *testing.T) {
	// Создаем парсер с тестовыми временами операций
	opTimes := OperationTimes{
		Addition:       100,
		Subtraction:    100,
		Multiplication: 200,
		Division:       200,
	}
	parser := NewParser(opTimes)

	tests := []struct {
		name     string
		expr     string
		wantErr  bool
		tasksLen int // Ожидаемое количество задач
	}{
		{
			name:     "Простое выражение",
			expr:     "2+2",
			wantErr:  false,
			tasksLen: 1,
		},
		{
			name:     "Выражение с приоритетом операций",
			expr:     "2+2*2",
			wantErr:  false,
			tasksLen: 2,
		},
		{
			name:     "Выражение со скобками",
			expr:     "(2+2)*2",
			wantErr:  false,
			tasksLen: 2,
		},
		{
			name:     "Сложное выражение",
			expr:     "2*(3+4)/(5-2)",
			wantErr:  false,
			tasksLen: 4,
		},
		{
			name:     "Пустое выражение",
			expr:     "",
			wantErr:  true,
			tasksLen: 0,
		},
		{
			name:     "Некорректное выражение",
			expr:     "2++2",
			wantErr:  true,
			tasksLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := parser.ParseExpression(tt.expr)

			// Проверяем наличие ошибки
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Если ошибка ожидается, не проверяем дальше
			if tt.wantErr {
				return
			}

			// Проверяем количество задач
			if len(tasks) != tt.tasksLen {
				t.Errorf("ParseExpression() tasksLen = %v, want %v", len(tasks), tt.tasksLen)
			}

			// Проверяем, что все задачи имеют правильные поля
			for i, task := range tasks {
				if task.Operation == "" {
					t.Errorf("Task %d has empty operation", i)
				}
				if task.Arg1 == "" {
					t.Errorf("Task %d has empty arg1", i)
				}
				if task.Arg2 == "" {
					t.Errorf("Task %d has empty arg2", i)
				}
			}
		})
	}
}

func TestStorage_AddExpression(t *testing.T) {
	storage := NewStorage()

	// Добавляем выражение
	id, err := storage.AddExpression("2+2")
	if err != nil {
		t.Errorf("AddExpression() error = %v", err)
		return
	}

	// Проверяем, что ID присвоен
	if id <= 0 {
		t.Errorf("AddExpression() id = %v, want > 0", id)
	}

	// Получаем выражение
	expr, err := storage.GetExpression(id)
	if err != nil {
		t.Errorf("GetExpression() error = %v", err)
		return
	}

	// Проверяем поля выражения
	if expr.ID != id {
		t.Errorf("Expression.ID = %v, want %v", expr.ID, id)
	}
	if expr.RawExpr != "2+2" {
		t.Errorf("Expression.RawExpr = %v, want %v", expr.RawExpr, "2+2")
	}
	if expr.Status != models.StatusPending {
		t.Errorf("Expression.Status = %v, want %v", expr.Status, models.StatusPending)
	}
}
