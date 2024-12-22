package calculator

import (
	"fmt"
	"strconv"
)

// TokenType представляет тип токена.
type TokenType string

// Поддерживаемые типы токенов.
const (
	TokenPlus     TokenType = "+"
	TokenMinus    TokenType = "-"
	TokenMultiply TokenType = "*"
	TokenDivide   TokenType = "/"
	TokenLParen   TokenType = "("
	TokenRParen   TokenType = ")"
	TokenNumber   TokenType = "NUMBER"
)

// Token представляет токен в выражении.
type Token struct {
	Type TokenType // Тип токена
	Name string    // Строковое представление токена
}

// Calc вычисляет результат арифметического выражения.
// Возвращает результат вычисления и ошибку, если она возникла.
func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)

	result, pos, err := parseExpression(tokens, 0)
	if err != nil {
		return 0, err
	}

	if pos != len(tokens) {
		return 0, fmt.Errorf("unexpected token: %s", tokens[pos].Name)
	}

	return result, nil
}

// tokenize разбивает строку на токены.
func tokenize(expression string) []Token {
	var tokens []Token
	var currentToken string

	for _, char := range expression {
		switch char {
		case '+', '-', '*', '/', '(', ')':
			if currentToken != "" {
				tokens = append(tokens, Token{Type: TokenNumber, Name: currentToken})
				currentToken = ""
			}
			tokens = append(tokens, Token{Type: TokenType(char), Name: string(char)})
		case ' ', '\t':
			if currentToken != "" {
				tokens = append(tokens, Token{Type: TokenNumber, Name: currentToken})
				currentToken = ""
			}
		default:
			currentToken += string(char)
		}
	}

	if currentToken != "" {
		tokens = append(tokens, Token{Type: TokenNumber, Name: currentToken})
	}

	return tokens
}

// parseExpression обрабатывает сложение и вычитание.
func parseExpression(tokens []Token, pos int) (float64, int, error) {
	left, pos, err := parseTerm(tokens, pos)
	if err != nil {
		return 0, pos, err
	}

	for pos < len(tokens) {
		switch tokens[pos].Type {
		case TokenPlus:
			right, newPos, err := parseTerm(tokens, pos+1)
			if err != nil {
				return 0, newPos, err
			}
			left += right
			pos = newPos
		case TokenMinus:
			right, newPos, err := parseTerm(tokens, pos+1)
			if err != nil {
				return 0, newPos, err
			}
			left -= right
			pos = newPos
		default:
			return left, pos, nil
		}
	}

	return left, pos, nil
}

// parseTerm обрабатывает умножение и деление.
func parseTerm(tokens []Token, pos int) (float64, int, error) {
	left, pos, err := parseFactor(tokens, pos)
	if err != nil {
		return 0, pos, err
	}

	for pos < len(tokens) {
		switch tokens[pos].Type {
		case TokenMultiply:
			right, newPos, err := parseFactor(tokens, pos+1)
			if err != nil {
				return 0, newPos, err
			}
			left *= right
			pos = newPos
		case TokenDivide:
			right, newPos, err := parseFactor(tokens, pos+1)
			if err != nil {
				return 0, newPos, err
			}
			if right == 0 {
				return 0, newPos, fmt.Errorf("division by zero")
			}
			left /= right
			pos = newPos
		default:
			return left, pos, nil
		}
	}

	return left, pos, nil
}

// parseFactor обрабатывает числа, унарные операторы и скобки.
// Возвращает значение, позицию после обработки и ошибку, если она возникла.
func parseFactor(tokens []Token, pos int) (float64, int, error) {
	if pos >= len(tokens) {
		return 0, pos, fmt.Errorf("insufficient tokens")
	}

	// Обработка унарных операторов
	if tokens[pos].Type == TokenPlus || tokens[pos].Type == TokenMinus {
		operator := tokens[pos].Type
		value, newPos, err := parseFactor(tokens, pos+1)
		if err != nil {
			return 0, newPos, err
		}
		if operator == TokenMinus {
			return -value, newPos, nil
		}
		return value, newPos, nil
	}

	// Обработка скобок
	if tokens[pos].Type == TokenLParen {
		result, newPos, err := parseExpression(tokens, pos+1)
		if err != nil {
			return 0, newPos, err
		}
		if newPos >= len(tokens) || tokens[newPos].Type != TokenRParen {
			return 0, newPos, fmt.Errorf("missing closing parenthesis")
		}
		return result, newPos + 1, nil
	}

	// Обработка чисел
	if tokens[pos].Type != TokenNumber {
		return 0, pos, fmt.Errorf("invalid token: %s", tokens[pos].Name)
	}

	num, err := strconv.ParseFloat(tokens[pos].Name, 64)
	if err != nil {
		return 0, pos, fmt.Errorf("invalid number: %s", tokens[pos].Name)
	}

	return num, pos + 1, nil
}
