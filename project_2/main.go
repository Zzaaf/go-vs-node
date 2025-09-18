package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Структура для хранения результата парсинга
type ParseResult struct {
	URL     string        // URL страницы
	Title   string        // Заголовок страницы
	Error   error         // Ошибка при парсинге
	Elapsed time.Duration // Время выполнения
}

// Функция для парсинга заголовка страницы
func parsePageTitle(url string) ParseResult {
	start := time.Now()

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Выполняем HTTP запрос
	resp, err := client.Get(url)
	if err != nil {
		return ParseResult{
			URL:     url,
			Title:   "",
			Error:   err,
			Elapsed: time.Since(start),
		}
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ParseResult{
			URL:     url,
			Title:   "",
			Error:   err,
			Elapsed: time.Since(start),
		}
	}

	// Ищем заголовок в HTML
	titleRegex := regexp.MustCompile(`<title[^>]*>([^<]*)</title>`)
	matches := titleRegex.FindStringSubmatch(string(body))

	var title string
	if len(matches) > 1 {
		title = strings.TrimSpace(matches[1])
	} else {
		title = "Заголовок не найден"
	}

	return ParseResult{
		URL:     url,
		Title:   title,
		Error:   nil,
		Elapsed: time.Since(start),
	}
}

// Последовательная версия парсера (медленная)
func parseSequential(urls []string) []ParseResult {
	fmt.Println("🐌 Запуск последовательного парсинга...")
	start := time.Now()

	var results []ParseResult
	for _, url := range urls {
		fmt.Printf("   Парсинг: %s\n", url)
		result := parsePageTitle(url)
		results = append(results, result)
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ Последовательный парсинг завершен за: %v\n\n", elapsed)

	return results
}

// Параллельная версия парсера (быстрая)
func parseParallel(urls []string) []ParseResult {
	fmt.Println("🚀 Запуск параллельного парсинга...")
	start := time.Now()

	// Канал для сбора результатов
	resultsChan := make(chan ParseResult, len(urls))

	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// Запускаем горутину для каждого URL
	for _, url := range urls {
		wg.Add(1) // Увеличиваем счетчик ожидаемых горутин

		go func(url string) {
			defer wg.Done() // Уменьшаем счетчик при завершении горутины

			fmt.Printf("   Парсинг: %s (в горутине)\n", url)
			result := parsePageTitle(url)
			resultsChan <- result // Отправляем результат в канал
		}(url)
	}

	// Запускаем горутину для закрытия канала после завершения всех рабочих горутин
	go func() {
		wg.Wait()          // Ждем завершения всех горутин
		close(resultsChan) // Закрываем канал
	}()

	// Собираем результаты из канала
	var results []ParseResult
	for result := range resultsChan {
		results = append(results, result)
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ Параллельный парсинг завершен за: %v\n\n", elapsed)

	return results
}

// Функция для вывода результатов
func printResults(results []ParseResult, title string) {
	fmt.Printf("📊 %s:\n", title)
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("❌ %s: Ошибка - %v (время: %v)\n",
				result.URL, result.Error, result.Elapsed)
		} else {
			fmt.Printf("✅ %s: %s (время: %v)\n",
				result.URL, result.Title, result.Elapsed)
		}
	}
	fmt.Println()
}

func main() {
	fmt.Println("🎯 Демонстрация многопоточности в Go")
	fmt.Println("=====================================")
	fmt.Println()

	// Список URL для парсинга
	urls := []string{
		"https://httpbin.org/html",
		"https://httpbin.org/json",
		"https://httpbin.org/xml",
		"https://httpbin.org/robots.txt",
		"https://httpbin.org/user-agent",
	}

	fmt.Printf("📝 Парсим %d URL...\n\n", len(urls))

	// Запускаем последовательный парсинг
	sequentialResults := parseSequential(urls)

	// Запускаем параллельный парсинг
	parallelResults := parseParallel(urls)

	// Выводим результаты
	printResults(sequentialResults, "Результаты последовательного парсинга")
	printResults(parallelResults, "Результаты параллельного парсинга")

	// Сравниваем производительность
	fmt.Println("📈 Сравнение производительности:")
	fmt.Println("=================================")

	// Вычисляем общее время для каждого метода
	var sequentialTime, parallelTime time.Duration
	for _, result := range sequentialResults {
		sequentialTime += result.Elapsed
	}
	for _, result := range parallelResults {
		parallelTime += result.Elapsed
	}

	fmt.Printf("Последовательный: %v\n", sequentialTime)
	fmt.Printf("Параллельный: %v\n", parallelTime)

	if parallelTime > 0 {
		speedup := float64(sequentialTime) / float64(parallelTime)
		fmt.Printf("Ускорение: %.2fx\n", speedup)
	}

	fmt.Println()
	fmt.Println("🎉 Ключевые концепции Go:")
	fmt.Println("• go func() - запуск горутины")
	fmt.Println("• sync.WaitGroup - ожидание завершения горутин")
	fmt.Println("• chan - каналы для безопасной передачи данных")
	fmt.Println("• Горутины очень дешевые - можно создавать тысячи!")
}
