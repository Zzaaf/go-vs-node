package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Структура для задачи
type Task struct {
	ID          int       `json:"id"`          // Уникальный идентификатор
	Title       string    `json:"title"`       // Название задачи
	Description string    `json:"description"` // Описание задачи
	Status      string    `json:"status"`      // Статус: "pending", "in_progress", "completed"
	Priority    string    `json:"priority"`    // Приоритет: "low", "medium", "high"
	CreatedAt   time.Time `json:"created_at"`  // Время создания
	UpdatedAt   time.Time `json:"updated_at"`  // Время последнего обновления
}

// Структура для создания новой задачи
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// Структура для обновления задачи
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *string `json:"priority,omitempty"`
}

// Структура для ответа API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// In-memory хранилище задач
type TaskStore struct {
	tasks  map[int]Task
	nextID int
	mutex  sync.RWMutex
}

// Создаем новое хранилище
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

// Добавляем новую задачу
func (ts *TaskStore) CreateTask(req CreateTaskRequest) Task {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Валидация приоритета
	if req.Priority == "" {
		req.Priority = "medium"
	}
	if req.Priority != "low" && req.Priority != "medium" && req.Priority != "high" {
		req.Priority = "medium"
	}

	now := time.Now()
	task := Task{
		ID:          ts.nextID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
		Priority:    req.Priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ts.tasks[ts.nextID] = task
	ts.nextID++

	return task
}

// Получаем все задачи
func (ts *TaskStore) GetAllTasks() []Task {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	tasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// Получаем задачу по ID
func (ts *TaskStore) GetTaskByID(id int) (Task, bool) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	task, exists := ts.tasks[id]
	return task, exists
}

// Обновляем задачу
func (ts *TaskStore) UpdateTask(id int, req UpdateTaskRequest) (Task, bool) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return Task{}, false
	}

	// Обновляем поля, если они переданы
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		// Валидация статуса
		if *req.Status == "pending" || *req.Status == "in_progress" || *req.Status == "completed" {
			task.Status = *req.Status
		}
	}
	if req.Priority != nil {
		// Валидация приоритета
		if *req.Priority == "low" || *req.Priority == "medium" || *req.Priority == "high" {
			task.Priority = *req.Priority
		}
	}

	task.UpdatedAt = time.Now()
	ts.tasks[id] = task

	return task, true
}

// Удаляем задачу
func (ts *TaskStore) DeleteTask(id int) bool {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	_, exists := ts.tasks[id]
	if exists {
		delete(ts.tasks, id)
	}
	return exists
}

// Глобальное хранилище задач
var taskStore = NewTaskStore()

// Middleware для логирования запросов
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// Middleware для CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Обработчик для получения всех задач
func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tasks := taskStore.GetAllTasks()

	response := APIResponse{
		Success: true,
		Message: "Задачи успешно получены",
		Data:    tasks,
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик для получения задачи по ID
func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(path)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	task, exists := taskStore.GetTaskByID(id)
	if !exists {
		response := APIResponse{
			Success: false,
			Error:   "Задача не найдена",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Задача найдена",
		Data:    task,
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик для создания новой задачи
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Неверный формат JSON",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Валидация обязательных полей
	if req.Title == "" {
		response := APIResponse{
			Success: false,
			Error:   "Название задачи обязательно",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	task := taskStore.CreateTask(req)

	response := APIResponse{
		Success: true,
		Message: "Задача успешно создана",
		Data:    task,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Обработчик для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(path)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Неверный формат JSON",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	task, exists := taskStore.UpdateTask(id, req)
	if !exists {
		response := APIResponse{
			Success: false,
			Error:   "Задача не найдена",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Задача успешно обновлена",
		Data:    task,
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(path)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	exists := taskStore.DeleteTask(id)
	if !exists {
		response := APIResponse{
			Success: false,
			Error:   "Задача не найдена",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Задача успешно удалена",
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик для корневого маршрута
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Добро пожаловать в API управления задачами!",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"endpoints": map[string]string{
				"GET    /tasks":      "Получить все задачи",
				"GET    /tasks/{id}": "Получить задачу по ID",
				"POST   /tasks":      "Создать новую задачу",
				"PUT    /tasks/{id}": "Обновить задачу",
				"DELETE /tasks/{id}": "Удалить задачу",
			},
		},
	}

	json.NewEncoder(w).Encode(response)
}

// Функция для настройки маршрутов
func setupRoutes() {
	// Применяем middleware ко всем обработчикам
	http.HandleFunc("/", corsMiddleware(loggingMiddleware(homeHandler)))
	http.HandleFunc("/tasks", corsMiddleware(loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/tasks" {
			getAllTasksHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})))
	http.HandleFunc("/tasks/", corsMiddleware(loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTaskByIDHandler(w, r)
		} else if r.Method == http.MethodPut {
			updateTaskHandler(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTaskHandler(w, r)
		} else {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})))
	http.HandleFunc("/create", corsMiddleware(loggingMiddleware(createTaskHandler)))
}

func main() {
	// Настраиваем маршруты
	setupRoutes()

	// Добавляем несколько тестовых задач
	taskStore.CreateTask(CreateTaskRequest{
		Title:       "Изучить Go",
		Description: "Изучить основы языка программирования Go",
		Priority:    "high",
	})

	taskStore.CreateTask(CreateTaskRequest{
		Title:       "Написать API",
		Description: "Создать REST API для управления задачами",
		Priority:    "medium",
	})

	// Запускаем сервер
	port := ":8080"
	fmt.Printf("🚀 REST API сервер запущен на http://localhost%s\n", port)
	fmt.Printf("📊 Доступные эндпоинты:\n")
	fmt.Printf("   GET    / - информация об API\n")
	fmt.Printf("   GET    /tasks - получить все задачи\n")
	fmt.Printf("   GET    /tasks/{id} - получить задачу по ID\n")
	fmt.Printf("   POST   /create - создать новую задачу\n")
	fmt.Printf("   PUT    /tasks/{id} - обновить задачу\n")
	fmt.Printf("   DELETE /tasks/{id} - удалить задачу\n")
	fmt.Printf("\n✅ Используется только стандартная библиотека Go!\n")
	fmt.Printf("🔧 Для тестирования используйте Postman или curl\n")

	log.Fatal(http.ListenAndServe(port, nil))
}
