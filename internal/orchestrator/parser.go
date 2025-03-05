package orchestrator

import (
	"fmt"
	"github.com/mpkelevra23/arithmetic-web-service/internal/models"
	"strconv"
	"strings"
)

// OperationTimes хранит время выполнения операций в миллисекундах
type OperationTimes struct {
	Addition       int
	Subtraction    int
	Multiplication int
	Division       int
}

// Parser представляет парсер арифметических выражений
type Parser struct {
	opTimes OperationTimes
}

// NewParser создает новый парсер
func NewParser(opTimes OperationTimes) *Parser {
	return &Parser{
		opTimes: opTimes,
	}
}

// Token представляет токен в выражении
type Token struct {
	Type  string // "NUMBER", "OPERATOR", "LPAREN", "RPAREN"
	Value string
}

// Node представляет узел в дереве выражения
type Node struct {
	Type   string // "NUMBER", "OPERATION"
	Value  string
	Left   *Node
	Right  *Node
	TaskID int
}

// ParseExpression разбирает выражение и создает задачи
func (p *Parser) ParseExpression(expr string) ([]models.Task, error) {
	// Удаляем пробелы
	expr = strings.ReplaceAll(expr, " ", "")

	if expr == "" {
		return nil, fmt.Errorf("пустое выражение")
	}

	// Разбиваем на токены
	tokens, err := p.tokenize(expr)
	if err != nil {
		return nil, err
	}

	// Строим дерево выражения
	root, remainingTokens, err := p.parseExpression(tokens, 0)
	if err != nil {
		return nil, err
	}

	if len(remainingTokens) > 0 {
		return nil, fmt.Errorf("некорректное выражение: лишние символы")
	}

	// Преобразуем дерево в задачи
	tasks := make([]models.Task, 0)
	_, err = p.buildTasks(root, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// tokenize разбивает строку на токены
func (p *Parser) tokenize(expr string) ([]Token, error) {
	tokens := make([]Token, 0)
	i := 0

	for i < len(expr) {
		char := expr[i]

		switch {
		case char >= '0' && char <= '9' || char == '.':
			// Число
			start := i
			hasDot := char == '.'

			i++
			for i < len(expr) && (expr[i] >= '0' && expr[i] <= '9' || expr[i] == '.') {
				if expr[i] == '.' {
					if hasDot {
						return nil, fmt.Errorf("некорректное число: две десятичные точки")
					}
					hasDot = true
				}
				i++
			}

			tokens = append(tokens, Token{Type: "NUMBER", Value: expr[start:i]})
		case char == '+' || char == '-' || char == '*' || char == '/':
			// Оператор
			tokens = append(tokens, Token{Type: "OPERATOR", Value: string(char)})
			i++
		case char == '(':
			// Левая скобка
			tokens = append(tokens, Token{Type: "LPAREN", Value: "("})
			i++
		case char == ')':
			// Правая скобка
			tokens = append(tokens, Token{Type: "RPAREN", Value: ")"})
			i++
		default:
			return nil, fmt.Errorf("некорректный символ: %c", char)
		}
	}

	return tokens, nil
}

// parseExpression разбирает выражение (сложение и вычитание)
func (p *Parser) parseExpression(tokens []Token, pos int) (*Node, []Token, error) {
	// Разбираем первый терм
	left, tokens, err := p.parseTerm(tokens, pos)
	if err != nil {
		return nil, tokens, err
	}

	// Пока есть операторы + или -
	for len(tokens) > 0 && (tokens[0].Value == "+" || tokens[0].Value == "-") {
		operator := tokens[0].Value
		tokens = tokens[1:] // Удаляем оператор

		// Разбираем следующий терм
		right, newTokens, err := p.parseTerm(tokens, 0)
		if err != nil {
			return nil, newTokens, err
		}

		// Создаем узел для операции
		left = &Node{
			Type:  "OPERATION",
			Value: operator,
			Left:  left,
			Right: right,
		}

		tokens = newTokens
	}

	return left, tokens, nil
}

// parseTerm разбирает терм (умножение и деление)
func (p *Parser) parseTerm(tokens []Token, pos int) (*Node, []Token, error) {
	// Разбираем первый фактор
	left, tokens, err := p.parseFactor(tokens, pos)
	if err != nil {
		return nil, tokens, err
	}

	// Пока есть операторы * или /
	for len(tokens) > 0 && (tokens[0].Value == "*" || tokens[0].Value == "/") {
		operator := tokens[0].Value
		tokens = tokens[1:] // Удаляем оператор

		// Разбираем следующий фактор
		right, newTokens, err := p.parseFactor(tokens, 0)
		if err != nil {
			return nil, newTokens, err
		}

		// Создаем узел для операции
		left = &Node{
			Type:  "OPERATION",
			Value: operator,
			Left:  left,
			Right: right,
		}

		tokens = newTokens
	}

	return left, tokens, nil
}

// parseFactor разбирает фактор (число или выражение в скобках)
func (p *Parser) parseFactor(tokens []Token, pos int) (*Node, []Token, error) {
	if len(tokens) == 0 {
		return nil, tokens, fmt.Errorf("неожиданный конец выражения")
	}

	token := tokens[0]
	tokens = tokens[1:] // Удаляем текущий токен

	switch token.Type {
	case "NUMBER":
		// Создаем узел-число
		return &Node{Type: "NUMBER", Value: token.Value}, tokens, nil
	case "LPAREN":
		// Разбираем выражение в скобках
		expr, newTokens, err := p.parseExpression(tokens, 0)
		if err != nil {
			return nil, newTokens, err
		}

		if len(newTokens) == 0 || newTokens[0].Type != "RPAREN" {
			return nil, newTokens, fmt.Errorf("ожидалась закрывающая скобка")
		}

		return expr, newTokens[1:], nil
	default:
		return nil, tokens, fmt.Errorf("неожиданный токен: %s", token.Value)
	}
}

// buildTasks преобразует дерево выражения в список задач
func (p *Parser) buildTasks(node *Node, tasks *[]models.Task, exprID int) (string, error) {
	if node.Type == "NUMBER" {
		// Для числа просто возвращаем его значение
		return node.Value, nil
	}

	// Рекурсивно обрабатываем левое и правое поддерево
	leftArg, err := p.buildTasks(node.Left, tasks, exprID)
	if err != nil {
		return "", err
	}

	rightArg, err := p.buildTasks(node.Right, tasks, exprID)
	if err != nil {
		return "", err
	}

	// Создаем задачу для текущей операции
	var operation models.Operation
	var operationTime int

	switch node.Value {
	case "+":
		operation = models.OperationAdd
		operationTime = p.opTimes.Addition
	case "-":
		operation = models.OperationSubtract
		operationTime = p.opTimes.Subtraction
	case "*":
		operation = models.OperationMultiply
		operationTime = p.opTimes.Multiplication
	case "/":
		operation = models.OperationDivide
		operationTime = p.opTimes.Division
	default:
		return "", fmt.Errorf("неизвестная операция: %s", node.Value)
	}

	// Создаем задачу
	task := models.Task{
		ID:            len(*tasks) + 1, // Временный ID
		ExpressionID:  exprID,
		Arg1:          leftArg,
		Arg2:          rightArg,
		Operation:     operation,
		OperationTime: operationTime,
		Dependencies:  make([]int, 0),
	}

	// Добавляем зависимости
	if strings.HasPrefix(leftArg, "res:") {
		taskID, err := strconv.Atoi(leftArg[4:])
		if err == nil {
			task.Dependencies = append(task.Dependencies, taskID)
		}
	}

	if strings.HasPrefix(rightArg, "res:") {
		taskID, err := strconv.Atoi(rightArg[4:])
		if err == nil {
			task.Dependencies = append(task.Dependencies, taskID)
		}
	}

	// Добавляем задачу в список
	*tasks = append(*tasks, task)
	taskID := len(*tasks)
	node.TaskID = taskID

	// Возвращаем ссылку на результат
	return fmt.Sprintf("res:%d", taskID), nil
}
