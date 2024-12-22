package errors

import (
	"encoding/json"
	"net/http"
)

// APIError представляет структурированную ошибку для API.
type APIError struct {
	Message string `json:"error"`
}

// WriteErrorResponse отправляет ошибку в формате JSON клиенту.
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	apiErr := APIError{
		Message: message,
	}

	json.NewEncoder(w).Encode(apiErr)
}

// Предопределенные сообщения об ошибках.
const (
	ErrInvalidInput      = "Expression is not valid"
	ErrMissingField      = "Missing field: expression"
	ErrTooLongExpression = "Expression is too long"
	ErrInternalServer    = "Internal server error"
	ErrDivisionByZero    = "Division by zero"
	ErrInvalidExpression = "Invalid expression"
	ErrMalformedJSON     = "Malformed JSON"
	ErrUnsupportedMethod = "Unsupported HTTP method"
)
