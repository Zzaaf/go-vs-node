package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Структура для хранения информации о файле
type FileInfo struct {
	Path     string    // Путь к файлу
	Size     int64     // Размер файла
	Hash     string    // MD5 хеш файла
	Modified time.Time // Время модификации
}

// Структура для группировки дубликатов
type DuplicateGroup struct {
	Hash  string     // Хеш группы
	Size  int64      // Размер файлов
	Files []FileInfo // Список файлов с одинаковым хешем
	Count int        // Количество дубликатов
}

// Функция для вычисления MD5 хеша файла
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Функция для сканирования директории и поиска файлов
func scanDirectory(dirPath string, extensions []string) ([]FileInfo, error) {
	var files []FileInfo

	fmt.Printf("🔍 Сканирование директории: %s\n", dirPath)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Фильтруем по расширениям, если они указаны
		if len(extensions) > 0 {
			ext := strings.ToLower(filepath.Ext(path))
			found := false
			for _, allowedExt := range extensions {
				if ext == strings.ToLower(allowedExt) {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		// Вычисляем хеш файла
		hash, err := calculateMD5(path)
		if err != nil {
			fmt.Printf("⚠️  Ошибка при чтении файла %s: %v\n", path, err)
			return nil
		}

		files = append(files, FileInfo{
			Path:     path,
			Size:     info.Size(),
			Hash:     hash,
			Modified: info.ModTime(),
		})

		return nil
	})

	return files, err
}

// Функция для поиска дубликатов
func findDuplicates(files []FileInfo) []DuplicateGroup {
	fmt.Printf("🔍 Поиск дубликатов среди %d файлов...\n", len(files))

	// Группируем файлы по хешу
	hashGroups := make(map[string][]FileInfo)
	for _, file := range files {
		hashGroups[file.Hash] = append(hashGroups[file.Hash], file)
	}

	// Создаем группы дубликатов
	var duplicates []DuplicateGroup
	for hash, groupFiles := range hashGroups {
		if len(groupFiles) > 1 {
			duplicates = append(duplicates, DuplicateGroup{
				Hash:  hash,
				Size:  groupFiles[0].Size,
				Files: groupFiles,
				Count: len(groupFiles),
			})
		}
	}

	// Сортируем по размеру (самые большие дубликаты сначала)
	sort.Slice(duplicates, func(i, j int) bool {
		return duplicates[i].Size > duplicates[j].Size
	})

	return duplicates
}

// Функция для вывода результатов
func printResults(duplicates []DuplicateGroup, showDetails bool) {
	if len(duplicates) == 0 {
		fmt.Println("✅ Дубликаты не найдены!")
		return
	}

	fmt.Printf("🎯 Найдено %d групп дубликатов:\n\n", len(duplicates))

	totalWastedSpace := int64(0)

	for i, group := range duplicates {
		fmt.Printf("📁 Группа %d (хеш: %s, размер: %s):\n",
			i+1, group.Hash[:8]+"...", formatSize(group.Size))

		// Сортируем файлы по времени модификации (самый старый первый)
		sort.Slice(group.Files, func(a, b int) bool {
			return group.Files[a].Modified.Before(group.Files[b].Modified)
		})

		for j, file := range group.Files {
			status := "🔄"
			if j == 0 {
				status = "✅" // Самый старый файл - оставляем
			} else {
				status = "🗑️ " // Остальные - можно удалить
			}

			fmt.Printf("  %s %s (%s, %s)\n",
				status,
				file.Path,
				formatSize(file.Size),
				file.Modified.Format("2006-01-02 15:04:05"))
		}

		// Вычисляем потраченное место (все файлы кроме одного)
		wastedSpace := group.Size * int64(group.Count-1)
		totalWastedSpace += wastedSpace

		fmt.Printf("  💾 Можно освободить: %s\n\n", formatSize(wastedSpace))
	}

	fmt.Printf("📊 Общая статистика:\n")
	fmt.Printf("   Всего групп дубликатов: %d\n", len(duplicates))
	fmt.Printf("   Можно освободить места: %s\n", formatSize(totalWastedSpace))
}

// Функция для форматирования размера файла
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Функция для вывода справки
func printHelp() {
	fmt.Println("🔍 Поисковик дубликатов файлов")
	fmt.Println("==============================")
	fmt.Println()
	fmt.Println("Использование:")
	fmt.Printf("  %s [опции] <директория>\n", os.Args[0])
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  -ext <расширения>  Фильтр по расширениям файлов (через запятую)")
	fmt.Println("  -details           Показать подробную информацию")
	fmt.Println("  -help              Показать эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Printf("  %s /path/to/folder\n", os.Args[0])
	fmt.Printf("  %s -ext .jpg,.png,.gif /path/to/images\n", os.Args[0])
	fmt.Printf("  %s -details /path/to/folder\n", os.Args[0])
	fmt.Println()
	fmt.Println("Сборка для разных платформ:")
	fmt.Println("  Windows:  GOOS=windows GOARCH=amd64 go build -o finder.exe")
	fmt.Println("  Linux:    GOOS=linux GOARCH=amd64 go build -o finder_linux")
	fmt.Println("  macOS:    GOOS=darwin GOARCH=arm64 go build -o finder_mac")
}

func main() {
	// Парсим аргументы командной строки
	var extensions string
	var showDetails bool
	var showHelp bool

	flag.StringVar(&extensions, "ext", "", "Фильтр по расширениям файлов (через запятую)")
	flag.BoolVar(&showDetails, "details", false, "Показать подробную информацию")
	flag.BoolVar(&showHelp, "help", false, "Показать справку")
	flag.Parse()

	// Показываем справку
	if showHelp {
		printHelp()
		return
	}

	// Получаем директорию для сканирования
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("❌ Ошибка: не указана директория для сканирования")
		fmt.Println("Используйте -help для получения справки")
		os.Exit(1)
	}

	dirPath := args[0]

	// Проверяем существование директории
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("❌ Ошибка: директория %s не существует\n", dirPath)
		os.Exit(1)
	}

	// Парсим расширения
	var extList []string
	if extensions != "" {
		extList = strings.Split(extensions, ",")
		for i, ext := range extList {
			extList[i] = strings.TrimSpace(ext)
			if !strings.HasPrefix(extList[i], ".") {
				extList[i] = "." + extList[i]
			}
		}
	}

	// Выводим информацию о программе
	fmt.Println("🔍 Поисковик дубликатов файлов")
	fmt.Println("==============================")
	fmt.Printf("Платформа: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Версия Go: %s\n", runtime.Version())
	fmt.Println()

	startTime := time.Now()

	// Сканируем директорию
	files, err := scanDirectory(dirPath, extList)
	if err != nil {
		fmt.Printf("❌ Ошибка при сканировании директории: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("📁 Файлы не найдены")
		return
	}

	fmt.Printf("📁 Найдено %d файлов\n\n", len(files))

	// Ищем дубликаты
	duplicates := findDuplicates(files)

	// Выводим результаты
	printResults(duplicates, showDetails)

	elapsed := time.Since(startTime)
	fmt.Printf("⏱️  Время выполнения: %v\n", elapsed)
}
