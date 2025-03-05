package orchestrator

import (
	"encoding/json"
	"fmt"
	"github.com/mpkelevra23/arithmetic-web-service/internal/models"
	"net/http"
	"strconv"
	"strings"
)

// Server представляет HTTP-сервер оркестратора
type Server struct {
	storage *Storage
	parser  *Parser
}

// NewServer создает новый сервер оркестратора
func NewServer(storage *Storage, parser *Parser) *Server {
	return &Server{
		storage: storage,
		parser:  parser,
	}
}

// SetupRoutes настраивает маршруты HTTP-сервера
func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// API для пользователей
	mux.HandleFunc("/api/v1/calculate", s.handleCalculate)
	mux.HandleFunc("/api/v1/expressions", s.handleGetExpressions)
	mux.HandleFunc("/api/v1/expressions/", s.handleGetExpression)

	// API для агентов
	mux.HandleFunc("/internal/task", s.handleTask)

	return mux
}

// handleCalculate обрабатывает запрос на добавление выражения
func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req models.ExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusUnprocessableEntity)
		return
	}

	if req.Expression == "" {
		http.Error(w, "Выражение не может быть пустым", http.StatusUnprocessableEntity)
		return
	}

	// Добавляем выражение в хранилище
	exprID, err := s.storage.AddExpression(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка добавления выражения: %v", err), http.StatusInternalServerError)
		return
	}

	// Разбираем выражение на задачи
	tasks, err := s.parser.ParseExpression(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка разбора выражения: %v", err), http.StatusUnprocessableEntity)
		return
	}

	// Добавляем задачи для выражения
	if err := s.storage.AddTasks(exprID, tasks); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка добавления задач: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	resp := models.ExpressionResponse{ID: exprID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// handleGetExpressions обрабатывает запрос на получение всех выражений
func (s *Server) handleGetExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	expressions := s.storage.GetAllExpressions()

	resp := models.ExpressionsResponse{Expressions: expressions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleGetExpression обрабатывает запрос на получение выражения по ID
func (s *Server) handleGetExpression(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	expr, err := s.storage.GetExpression(id)
	if err != nil {
		http.Error(w, "Выражение не найдено", http.StatusNotFound)
		return
	}

	resp := models.ExpressionDetailResponse{Expression: expr}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleTask обрабатывает запросы агентов
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Получение задачи
		task, err := s.storage.GetReadyTask()
		if err != nil {
			http.Error(w, "Нет доступных задач", http.StatusNotFound)
			return
		}

		resp := models.TaskResponse{Task: task}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodPost:
		// Получение результата задачи
		var req models.TaskResultRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Некорректный JSON", http.StatusUnprocessableEntity)
			return
		}

		if err := s.storage.UpdateTaskResult(req.ID, req.Result, req.Error); err != nil {
			if strings.Contains(err.Error(), "не найдена") {
				http.Error(w, "Задача не найдена", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
