package tests

import (
	"testing"

	"github.com/mpkelevra23/arithmetic-web-service/internal/calculator"
)

// TestCalc проверяет корректность работы функции Calc(expression string) (float64, error).
func TestCalc(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected float64
		err      bool
	}{
		{
			name:     "Простое сложение",
			expr:     "1 + 2",
			expected: 3,
			err:      false,
		},
		{
			name:     "Простое вычитание",
			expr:     "5 - 3",
			expected: 2,
			err:      false,
		},
		{
			name:     "Умножение и деление",
			expr:     "4 * 2 / 2",
			expected: 4,
			err:      false,
		},
		{
			name:     "Скобки с приоритетом",
			expr:     "(1 + 2) * 3",
			expected: 9,
			err:      false,
		},
		{
			name:     "Сложное выражение",
			expr:     "3 + 4 * 2 / (1 - 5) * 2 + 3",
			expected: 2,
			err:      false,
		},
		{
			name:     "Вещественные числа",
			expr:     "3.5 + 2.5",
			expected: 6,
			err:      false,
		},
		{
			name:     "Отрицательные числа",
			expr:     "-2 + 3",
			expected: 1,
			err:      false,
		},
		{
			name:     "Деление на ноль",
			expr:     "10 / (5 - 5)",
			expected: 0, // Результат не важен, так как ожидается ошибка
			err:      true,
		},
		{
			name:     "Отсутствие закрывающей скобки",
			expr:     "(1 + 2 * 3",
			expected: 0, // Результат не важен, так как ожидается ошибка
			err:      true,
		},
		{
			name:     "Неверный токен",
			expr:     "2 + a",
			expected: 0, // Результат не важен, так как ожидается ошибка
			err:      true,
		},
		{
			name:     "Пустая строка",
			expr:     "",
			expected: 0,    // Результат не важен, так как может быть ошибка или 0
			err:      true, // Предполагаем, что пустая строка невалидна
		},
		{
			name:     "Только число",
			expr:     "42",
			expected: 42,
			err:      false,
		},
		{
			name:     "Многоступенчатые скобки",
			expr:     "((2 + 3) * (4 - 1)) / 5",
			expected: 3, // (5 * 3)/5 = 3
			err:      false,
		},
		{
			name:     "Выражение с пробелами и табуляциями",
			expr:     "  7 \t* ( 8 + 2 ) ",
			expected: 70, // 7 * 10 = 70
			err:      false,
		},
		{
			name:     "Выражение с несколькими операторами подряд",
			expr:     "1 + 2 - 3 + 4",
			expected: 4, // 1+2=3; 3-3=0; 0+4=4
			err:      false,
		},
		{
			name:     "Выражение с отрицательным результатом",
			expr:     "2 - 5",
			expected: -3,
			err:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.Calc(tt.expr)
			if tt.err {
				// Проверяем наличие ожидаемой ошибки
				if err == nil {
					t.Errorf("Calc(%q) = %v, ожидается ошибка", tt.expr, result)
				}
			} else {
				// Проверяем отсутствие ошибки и соответствие результата
				if err != nil {
					t.Errorf("Calc(%q) вернул ошибку: %v, ожидается %v", tt.expr, err, tt.expected)
				} else if result != tt.expected {
					t.Errorf("Calc(%q) = %v; ожидается %v", tt.expr, result, tt.expected)
				}
			}
		})
	}
}
