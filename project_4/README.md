# Project 4: От идеи к API за 15 минут - REST API микросервис

## Описание
Демонстрация того, насколько стандартная библиотека Go богата и самодостаточна для создания полноценных REST API без внешних фреймворков.

## Что изучаем
- **REST API** - создание веб-API с нуля
- **Стандартная библиотека** - только `net/http` и `encoding/json`
- **JSON** - работа с JSON данными
- **HTTP методы** - GET, POST, PUT, DELETE
- **Middleware** - логирование и CORS
- **Валидация** - проверка входных данных

## Функциональность
API для управления задачами (Task Management):
- ✅ Создание задач
- ✅ Получение всех задач
- ✅ Получение задачи по ID
- ✅ Обновление задач
- ✅ Удаление задач
- ✅ Валидация данных
- ✅ CORS поддержка
- ✅ Логирование запросов

## Запуск

### Запуск сервера
```bash
go run main.go
```

### Сборка
```bash
go build -o api-server main.go
./api-server
```

## API Эндпоинты

### 1. Информация об API
```http
GET /
```

**Ответ:**
```json
{
  "success": true,
  "message": "Добро пожаловать в API управления задачами!",
  "data": {
    "version": "1.0.0",
    "endpoints": {
      "GET    /tasks": "Получить все задачи",
      "GET    /tasks/{id}": "Получить задачу по ID",
      "POST   /tasks": "Создать новую задачу",
      "PUT    /tasks/{id}": "Обновить задачу",
      "DELETE /tasks/{id}": "Удалить задачу"
    }
  }
}
```

### 2. Получить все задачи
```http
GET /tasks
```

**Ответ:**
```json
{
  "success": true,
  "message": "Задачи успешно получены",
  "data": [
    {
      "id": 1,
      "title": "Изучить Go",
      "description": "Изучить основы языка программирования Go",
      "status": "pending",
      "priority": "high",
      "created_at": "2023-12-01T10:00:00Z",
      "updated_at": "2023-12-01T10:00:00Z"
    }
  ]
}
```

### 3. Создать новую задачу
```http
POST /create
Content-Type: application/json

{
  "title": "Новая задача",
  "description": "Описание задачи",
  "priority": "medium"
}
```

**Ответ:**
```json
{
  "success": true,
  "message": "Задача успешно создана",
  "data": {
    "id": 2,
    "title": "Новая задача",
    "description": "Описание задачи",
    "status": "pending",
    "priority": "medium",
    "created_at": "2023-12-01T10:30:00Z",
    "updated_at": "2023-12-01T10:30:00Z"
  }
}
```

### 4. Получить задачу по ID
```http
GET /tasks/{id}
```

### 5. Обновить задачу
```http
PUT /tasks/{id}
Content-Type: application/json

{
  "status": "in_progress",
  "priority": "high"
}
```

### 6. Удалить задачу
```http
DELETE /tasks/{id}
```

## Тестирование

### Автоматическое тестирование (Linux/macOS)
```bash
chmod +x test_api.sh
./test_api.sh
```

### Автоматическое тестирование (Windows)
```bash
test_api.bat
```

### Ручное тестирование с curl

#### Получить все задачи
```bash
curl http://localhost:8080/tasks
```

#### Создать задачу
```bash
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -d '{"title": "Тест", "description": "Описание", "priority": "high"}'
```

#### Обновить задачу
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

## Ключевые особенности кода

### 1. Структуры данных
```go
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Priority    string    `json:"priority"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 2. In-memory хранилище
```go
type TaskStore struct {
    tasks  map[int]Task
    nextID int
    mutex  sync.RWMutex  // Потокобезопасность
}
```

### 3. Middleware
```go
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        // ... настройки CORS
        next(w, r)
    }
}
```

### 4. JSON обработка
```go
// Декодирование JSON
var req CreateTaskRequest
json.NewDecoder(r.Body).Decode(&req)

// Кодирование JSON
json.NewEncoder(w).Encode(response)
```

### 5. Валидация данных
```go
if req.Title == "" {
    response := APIResponse{
        Success: false,
        Error:   "Название задачи обязательно",
    }
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(response)
    return
}
```

## Преимущества Go для API

### 1. Стандартная библиотека
- `net/http` - полноценный HTTP сервер
- `encoding/json` - работа с JSON
- `sync` - потокобезопасность
- Никаких внешних зависимостей!

### 2. Производительность
- Быстрый запуск
- Низкое потребление памяти
- Высокая пропускная способность

### 3. Простота
- Читаемый код
- Легкая отладка
- Простое тестирование

### 4. Безопасность
- Потокобезопасность из коробки
- Валидация данных
- Обработка ошибок

## Сравнение с другими языками

| Язык | Фреймворк | Зависимости | Сложность |
|------|-----------|-------------|-----------|
| **Go** | net/http | 0 | ⭐⭐ |
| Node.js | Express | 10+ | ⭐⭐⭐ |
| Python | Flask | 5+ | ⭐⭐⭐ |
| Java | Spring | 50+ | ⭐⭐⭐⭐ |

## Заключение
Go позволяет создавать полноценные REST API за считанные минуты, используя только стандартную библиотеку. Никаких фреймворков, никаких зависимостей - только чистый, быстрый и надежный код!
