package calculator

import (
	"fmt"     // Пакет fmt используется для форматированного ввода и вывода данных.
	"strconv" // Пакет strconv предоставляет функции для преобразования типов данных, например, из строк в числа.
)

// TokenType представляет тип токена. Это алиас для строки, используемый для определения различных типов токенов.
type TokenType string

// Константы, представляющие различные типы токенов, используемых в арифметическом выражении.
const (
	TokenPlus     TokenType = "+"      // Токен для оператора сложения.
	TokenMinus    TokenType = "-"      // Токен для оператора вычитания.
	TokenMultiply TokenType = "*"      // Токен для оператора умножения.
	TokenDivide   TokenType = "/"      // Токен для оператора деления.
	TokenLParen   TokenType = "("      // Токен для открывающей скобки.
	TokenRParen   TokenType = ")"      // Токен для закрывающей скобки.
	TokenNumber   TokenType = "NUMBER" // Токен для числовых значений.
)

// Token представляет собой токен с типом и именем. Используется для хранения информации о каждом токене в выражении.
type Token struct {
	Type TokenType // Тип токена (например, оператор или число).
	Name string    // Строковое представление токена (например, "+", "-", "NUMBER").
}

// Calc вычисляет арифметическое выражение, представленное строкой.
// Возвращает результат вычисления как float64 и ошибку, если она возникает.
// Метод начинается с большой буквы, что означает, что он является публичным и доступным из других пакетов.
func Calc(expression string) (float64, error) {
	// Разбиваем выражение на токены с помощью функции tokenize.
	tokens := tokenize(expression)

	// Начинаем парсинг выражения с позиции 0, используя функцию parseExpression.
	result, pos, err := parseExpression(tokens, 0)
	if err != nil {
		return 0, err // Возвращаем ошибку, если она возникла во время парсинга.
	}

	// Проверяем, что все токены были обработаны. Если нет, возвращаем ошибку.
	if pos != len(tokens) {
		return 0, fmt.Errorf("неожиданный токен: %s", tokens[pos].Name)
	}

	return result, nil // Возвращаем результат вычисления.
}

// tokenize разбивает строку арифметического выражения на слайс токенов.
// Каждый токен представляет собой оператор, скобку или число.
// Метод начинается с маленькой буквы, что означает, что он является приватным и доступным только внутри пакета.
func tokenize(expression string) []Token {
	var tokens []Token      // Инициализируем пустой слайс токенов.
	var currentToken string // Переменная для накопления символов числа.

	for _, char := range expression { // Проходим по каждому символу в выражении.
		switch char {
		case '+', '-', '*', '/', '(', ')': // Проверяем, является ли символ оператором или скобкой.
			if currentToken != "" { // Если есть накопленные символы числа, добавляем их как токен.
				tokens = append(tokens, Token{Type: TokenNumber, Name: currentToken})
				currentToken = "" // Сбрасываем накопитель для числа.
			}
			// Добавляем оператор или скобку как отдельный токен.
			tokens = append(tokens, Token{Type: TokenType(char), Name: string(char)})
		case ' ', '\t': // Если символ является пробелом или табуляцией.
			if currentToken != "" { // Добавляем накопленное число как токен.
				tokens = append(tokens, Token{Type: TokenNumber, Name: currentToken})
				currentToken = "" // Сбрасываем накопитель для числа.
			}
			// Пробелы игнорируются.
		default:
			// Добавляем символ к текущему числу.
			currentToken += string(char)
		}
	}

	// После обработки всех символов, если есть накопленное число, добавляем его как токен.
	if currentToken != "" {
		tokens = append(tokens, Token{Type: TokenNumber, Name: currentToken})
	}

	return tokens // Возвращаем слайс токенов.
}

// parseExpression рекурсивно обрабатывает выражение, начиная с позиции pos.
// Возвращает результат вычисления, новую позицию и ошибку, если она возникла.
func parseExpression(tokens []Token, pos int) (float64, int, error) {
	// Парсим первый термин выражения.
	left, pos, err := parseTerm(tokens, pos)
	if err != nil {
		return 0, pos, err // Возвращаем ошибку, если она возникла.
	}

	// Обрабатываем операторы + и -.
	for pos < len(tokens) {
		switch tokens[pos].Type {
		case TokenPlus: // Если текущий токен - оператор сложения.
			// Парсим следующий термин после оператора +.
			right, newPos, err := parseTerm(tokens, pos+1)
			if err != nil {
				return 0, newPos, err // Возвращаем ошибку, если она возникла.
			}
			left += right // Добавляем правый термин к левому.
			pos = newPos  // Обновляем текущую позицию.
		case TokenMinus: // Если текущий токен - оператор вычитания.
			// Парсим следующий термин после оператора -.
			right, newPos, err := parseTerm(tokens, pos+1)
			if err != nil {
				return 0, newPos, err // Возвращаем ошибку, если она возникла.
			}
			left -= right // Вычитаем правый термин из левого.
			pos = newPos  // Обновляем текущую позицию.
		default:
			// Если оператор не + или -, возвращаем текущий результат.
			return left, pos, nil
		}
	}

	return left, pos, nil // Возвращаем итоговый результат и позицию.
}

// parseTerm обрабатывает умножение и деление, начиная с позиции pos.
// Возвращает результат вычисления, новую позицию и ошибку, если она возникла.
func parseTerm(tokens []Token, pos int) (float64, int, error) {
	// Парсим первый фактор термина.
	left, pos, err := parseFactor(tokens, pos)
	if err != nil {
		return 0, pos, err // Возвращаем ошибку, если она возникла.
	}

	// Обрабатываем операторы * и /.
	for pos < len(tokens) {
		switch tokens[pos].Type {
		case TokenMultiply: // Если текущий токен - оператор умножения.
			// Парсим следующий фактор после оператора *.
			right, newPos, err := parseFactor(tokens, pos+1)
			if err != nil {
				return 0, newPos, err // Возвращаем ошибку, если она возникла.
			}
			left *= right // Умножаем левый фактор на правый.
			pos = newPos  // Обновляем текущую позицию.
		case TokenDivide: // Если текущий токен - оператор деления.
			// Парсим следующий фактор после оператора /.
			right, newPos, err := parseFactor(tokens, pos+1)
			if err != nil {
				return 0, newPos, err // Возвращаем ошибку, если она возникла.
			}
			// Проверяем деление на ноль.
			if right == 0 {
				return 0, newPos, fmt.Errorf("деление на ноль")
			}
			left /= right // Делим левый фактор на правый.
			pos = newPos  // Обновляем текущую позицию.
		default:
			// Если оператор не * или /, возвращаем текущий результат.
			return left, pos, nil
		}
	}

	return left, pos, nil // Возвращаем итоговый результат и позицию.
}

// parseFactor обрабатывает числа и выражения в скобках, начиная с позиции pos.
// Возвращает результат вычисления, новую позицию и ошибку, если она возникла.
func parseFactor(tokens []Token, pos int) (float64, int, error) {
	// Проверяем, что текущая позиция не выходит за пределы списка токенов.
	if pos >= len(tokens) {
		return 0, pos, fmt.Errorf("недостаточно токенов")
	}

	// Обработка унарных операторов + и -.
	if tokens[pos].Type == TokenPlus || tokens[pos].Type == TokenMinus {
		operator := tokens[pos].Type // Сохраняем оператор.
		// Парсим следующий фактор после оператора.
		value, newPos, err := parseFactor(tokens, pos+1)
		if err != nil {
			return 0, newPos, err // Возвращаем ошибку, если она возникла.
		}
		if operator == TokenMinus {
			return -value, newPos, nil // Инвертируем знак, если оператор -.
		}
		return value, newPos, nil // Возвращаем значение без изменений, если оператор +.
	}

	// Проверяем, является ли токен открывающей скобкой.
	if tokens[pos].Type == TokenLParen {
		// Парсим выражение внутри скобок.
		result, newPos, err := parseExpression(tokens, pos+1)
		if err != nil {
			return 0, newPos, err // Возвращаем ошибку, если она возникла.
		}
		// Проверяем наличие закрывающей скобки.
		if newPos >= len(tokens) || tokens[newPos].Type != TokenRParen {
			return 0, newPos, fmt.Errorf("отсутствует закрывающая скобка")
		}
		// Возвращаем результат выражения внутри скобок и обновленную позицию.
		return result, newPos + 1, nil
	}

	// Если токен не является скобкой, предполагаем, что это число.
	if tokens[pos].Type != TokenNumber {
		return 0, pos, fmt.Errorf("недопустимый токен: %s", tokens[pos].Name)
	}

	// Преобразуем строковое представление числа в тип float64.
	num, err := strconv.ParseFloat(tokens[pos].Name, 64)
	if err != nil {
		return 0, pos, fmt.Errorf("недопустимый токен: %s", tokens[pos].Name)
	}
	// Возвращаем число и обновленную позицию.
	return num, pos + 1, nil
}
