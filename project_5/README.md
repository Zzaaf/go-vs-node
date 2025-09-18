# Project 5: Go в неожиданном месте - WebAssembly приложение

## Описание
Демонстрация того, что на Go можно писать не только для бэкенда, но и компилировать код в WebAssembly для выполнения в браузере с высокой производительностью.

## Что изучаем
- **WebAssembly** - компиляция Go в WASM
- **Интерактивность** - взаимодействие Go с JavaScript
- **Анимация** - создание анимированных эффектов
- **Canvas API** - работа с 2D графикой
- **Производительность** - быстрые вычисления в браузере

## Функциональность
Интерактивная анимация частиц:
- 🎨 50 анимированных частиц разных цветов
- 💥 Клик для создания взрыва частиц
- 🔗 Соединения между близкими частицами
- 📊 Информация о системе
- 🎮 Интерактивные элементы управления

## Сборка

### Автоматическая сборка (Linux/macOS)
```bash
chmod +x build.sh
./build.sh
```

### Автоматическая сборка (Windows)
```bash
build.bat
```

### Ручная сборка
```bash
# 1. Компиляция Go в WebAssembly
GOOS=js GOARCH=wasm go build -o main.wasm main.go

# 2. Копирование JavaScript runtime
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

# 3. Запуск HTTP сервера
go run server.go
```

## Запуск

### Вариант 1: Go HTTP сервер
```bash
go run server.go
```
Откройте: http://localhost:8000

### Вариант 2: Python HTTP сервер
```bash
python -m http.server 8000
```
Откройте: http://localhost:8000

### Вариант 3: Node.js HTTP сервер
```bash
npx http-server -p 8000
```
Откройте: http://localhost:8000

## Использование

### Интерактивность
- **Клик по canvas** - создание взрыва частиц
- **Кнопка "Сбросить"** - сброс к начальному состоянию
- **Кнопка "Добавить частицы"** - добавление новых частиц
- **Кнопка "Информация"** - показ системной информации

### JavaScript API
Приложение экспортирует функции в JavaScript:
```javascript
// Получить информацию о системе
const info = getSystemInfo();
console.log(info);

// Сбросить частицы
resetParticles();

// Добавить новые частицы
addParticles(20);
```

## Ключевые особенности кода

### 1. WebAssembly компиляция
```bash
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

### 2. Взаимодействие с JavaScript
```go
// Получение DOM элементов
canvas := js.Global().Get("document").Call("getElementById", "canvas")
ctx := canvas.Call("getContext", "2d")

// Регистрация обработчиков событий
canvas.Call("addEventListener", "click", js.FuncOf(handleClick))

// Экспорт функций в JavaScript
js.Global().Set("getSystemInfo", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    return getSystemInfo()
}))
```

### 3. Анимация с requestAnimationFrame
```go
func animate() {
    updateParticles()
    drawParticles()
    animation = js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        animate()
        return nil
    }))
}
```

### 4. Canvas 2D API
```go
// Очистка canvas
ctx.Call("clearRect", 0, 0, width, height)

// Рисование круга
ctx.Set("fillStyle", particle.Color)
ctx.Call("beginPath")
ctx.Call("arc", particle.Position.X, particle.Position.Y, particle.Size, 0, 2*math.Pi)
ctx.Call("fill")
```

## Преимущества Go WebAssembly

### 1. Производительность
- Нативная скорость выполнения
- Эффективное использование памяти
- Быстрые математические вычисления

### 2. Простота разработки
- Знакомый синтаксис Go
- Статическая типизация
- Отличная поддержка в IDE

### 3. Кроссплатформенность
- Один код для всех браузеров
- Нет проблем с совместимостью
- Легкое тестирование

### 4. Безопасность
- Песочница WebAssembly
- Нет прямого доступа к DOM
- Контролируемое взаимодействие

## Сравнение с другими технологиями

| Технология | Производительность | Сложность | Размер |
|------------|-------------------|-----------|--------|
| **Go WASM** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ |
| JavaScript | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| TypeScript | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| Rust WASM | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |

## Применение

### 1. Игры
- 2D/3D игры в браузере
- Физические симуляции
- Интерактивные визуализации

### 2. Научные вычисления
- Математические расчеты
- Обработка данных
- Визуализация результатов

### 3. Графические редакторы
- Обработка изображений
- Векторная графика
- Анимация

### 4. Образование
- Интерактивные демонстрации
- Симуляторы
- Обучающие игры

## Отладка

### 1. Консоль браузера
```javascript
// Проверка загрузки WASM
console.log('WASM загружен:', typeof getSystemInfo === 'function');

// Получение информации о системе
console.log(getSystemInfo());
```

### 2. DevTools
- **Sources** - просмотр main.wasm
- **Network** - проверка загрузки файлов
- **Console** - логи и ошибки

### 3. Go отладка
```go
// Логирование в консоль браузера
js.Global().Get("console").Call("log", "Сообщение из Go")
```

## Заключение
Go WebAssembly открывает новые возможности для веб-разработки. Теперь можно использовать мощь Go для создания высокопроизводительных веб-приложений, игр и интерактивных демонстраций прямо в браузере!
