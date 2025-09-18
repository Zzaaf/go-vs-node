package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Структура для ответа сервера
type Response struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Note      string `json:"note,omitempty"`
}

// Функция для имитации долгой операции (например, запрос к базе данных)
func simulateLongOperation() string {
	// Имитируем долгую операцию - 10 секунд
	time.Sleep(10 * time.Second)
	return "Долгая операция завершена!"
}

// Обработчик для быстрого маршрута
func fastHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Message:   "Привет от Go сервера!",
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "success",
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик для медленного маршрута
func slowHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("[%s] Начало обработки медленного запроса\n", time.Now().Format(time.RFC3339))

	// Выполняем долгую операцию
	result := simulateLongOperation()

	fmt.Printf("[%s] Завершение обработки медленного запроса\n", time.Now().Format(time.RFC3339))

	response := Response{
		Message:   result,
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "success",
		Note:      "Этот запрос НЕ блокирует другие горутины!",
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик для 404 ошибок
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	response := Response{
		Message: "Страница не найдена",
		Status:  "error",
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	// Настраиваем маршруты
	http.HandleFunc("/", fastHandler)
	http.HandleFunc("/slow", slowHandler)
	http.HandleFunc("/*", notFoundHandler)

	// Запускаем сервер на порту 8080
	PORT := ":8080"
	fmt.Printf("🚀 Go сервер запущен на http://localhost%s\n", PORT)
	fmt.Printf("📊 Тестовые маршруты:\n")
	fmt.Printf("   GET / - быстрый ответ\n")
	fmt.Printf("   GET /slow - медленный ответ (10 сек)\n")
	fmt.Printf("\n✅ Преимущество: Горутины позволяют обрабатывать множество запросов параллельно!\n")

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(PORT, nil))
}
