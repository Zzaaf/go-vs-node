package main

import (
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

// Структура для точки
type Point struct {
	X float64
	Y float64
}

// Структура для частицы
type Particle struct {
	Position Point
	Velocity Point
	Color    string
	Size     float64
}

// Глобальные переменные
var (
	canvas    js.Value
	ctx       js.Value
	particles []Particle
	width     float64
	height    float64
	animation js.Value
)

// Функция для инициализации частиц
func initParticles() {
	particles = make([]Particle, 50)
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8"}

	for i := range particles {
		particles[i] = Particle{
			Position: Point{
				X: rand.Float64() * width,
				Y: rand.Float64() * height,
			},
			Velocity: Point{
				X: (rand.Float64() - 0.5) * 2,
				Y: (rand.Float64() - 0.5) * 2,
			},
			Color: colors[rand.Intn(len(colors))],
			Size:  rand.Float64()*3 + 1,
		}
	}
}

// Функция для обновления частиц
func updateParticles() {
	for i := range particles {
		// Обновляем позицию
		particles[i].Position.X += particles[i].Velocity.X
		particles[i].Position.Y += particles[i].Velocity.Y

		// Отскок от границ
		if particles[i].Position.X < 0 || particles[i].Position.X > width {
			particles[i].Velocity.X *= -1
		}
		if particles[i].Position.Y < 0 || particles[i].Position.Y > height {
			particles[i].Velocity.Y *= -1
		}

		// Ограничиваем позицию в пределах canvas
		particles[i].Position.X = math.Max(0, math.Min(width, particles[i].Position.X))
		particles[i].Position.Y = math.Max(0, math.Min(height, particles[i].Position.Y))
	}
}

// Функция для отрисовки частиц
func drawParticles() {
	// Очищаем canvas
	ctx.Call("clearRect", 0, 0, width, height)

	// Рисуем частицы
	for _, particle := range particles {
		ctx.Set("fillStyle", particle.Color)
		ctx.Call("beginPath")
		ctx.Call("arc", particle.Position.X, particle.Position.Y, particle.Size, 0, 2*math.Pi)
		ctx.Call("fill")
	}

	// Рисуем соединения между близкими частицами
	drawConnections()
}

// Функция для отрисовки соединений между частицами
func drawConnections() {
	ctx.Set("strokeStyle", "rgba(255, 255, 255, 0.1)")
	ctx.Set("lineWidth", 1)

	for i := 0; i < len(particles); i++ {
		for j := i + 1; j < len(particles); j++ {
			dx := particles[i].Position.X - particles[j].Position.X
			dy := particles[i].Position.Y - particles[j].Position.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			// Рисуем линию, если частицы близко
			if distance < 100 {
				opacity := 1 - distance/100
				ctx.Set("strokeStyle", js.ValueOf("rgba(255, 255, 255, "+js.ValueOf(opacity).String()+")"))
				ctx.Call("beginPath")
				ctx.Call("moveTo", particles[i].Position.X, particles[i].Position.Y)
				ctx.Call("lineTo", particles[j].Position.X, particles[j].Position.Y)
				ctx.Call("stroke")
			}
		}
	}
}

// Функция анимации
func animate() {
	updateParticles()
	drawParticles()
	animation = js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		animate()
		return nil
	}))
}

// Функция для обработки клика мыши
func handleClick(this js.Value, args []js.Value) interface{} {
	event := args[0]
	rect := canvas.Call("getBoundingClientRect")

	// Получаем координаты клика относительно canvas
	clickX := event.Get("clientX").Float() - rect.Get("left").Float()
	clickY := event.Get("clientY").Float() - rect.Get("top").Float()

	// Создаем взрыв частиц в точке клика
	createExplosion(clickX, clickY)

	return nil
}

// Функция для создания взрыва частиц
func createExplosion(x, y float64) {
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8"}

	// Добавляем 10 новых частиц в точке клика
	for i := 0; i < 10; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*5 + 2

		particle := Particle{
			Position: Point{X: x, Y: y},
			Velocity: Point{
				X: math.Cos(angle) * speed,
				Y: math.Sin(angle) * speed,
			},
			Color: colors[rand.Intn(len(colors))],
			Size:  rand.Float64()*4 + 2,
		}

		particles = append(particles, particle)
	}

	// Ограничиваем количество частиц
	if len(particles) > 200 {
		particles = particles[len(particles)-200:]
	}
}

// Функция для обработки изменения размера окна
func handleResize(this js.Value, args []js.Value) interface{} {
	width = js.Global().Get("innerWidth").Float()
	height = js.Global().Get("innerHeight").Float()

	canvas.Set("width", width)
	canvas.Set("height", height)

	return nil
}

// Функция для получения информации о системе
func getSystemInfo() map[string]interface{} {
	return map[string]interface{}{
		"userAgent": js.Global().Get("navigator").Get("userAgent").String(),
		"language":  js.Global().Get("navigator").Get("language").String(),
		"platform":  js.Global().Get("navigator").Get("platform").String(),
		"width":     width,
		"height":    height,
		"particles": len(particles),
	}
}

// Функция для экспорта в JavaScript
func registerCallbacks() {
	// Экспортируем функции в глобальный объект window
	js.Global().Set("getSystemInfo", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return getSystemInfo()
	}))

	js.Global().Set("resetParticles", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		initParticles()
		return nil
	}))

	js.Global().Set("addParticles", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		count := 10
		if len(args) > 0 {
			count = args[0].Int()
		}

		colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8"}

		for i := 0; i < count; i++ {
			particle := Particle{
				Position: Point{
					X: rand.Float64() * width,
					Y: rand.Float64() * height,
				},
				Velocity: Point{
					X: (rand.Float64() - 0.5) * 2,
					Y: (rand.Float64() - 0.5) * 2,
				},
				Color: colors[rand.Intn(len(colors))],
				Size:  rand.Float64()*3 + 1,
			}
			particles = append(particles, particle)
		}

		return nil
	}))
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Получаем canvas элемент
	canvas = js.Global().Get("document").Call("getElementById", "canvas")
	if canvas.IsUndefined() {
		js.Global().Get("console").Call("error", "Canvas элемент не найден!")
		return
	}

	// Получаем контекст 2D
	ctx = canvas.Call("getContext", "2d")

	// Получаем размеры окна
	width = js.Global().Get("innerWidth").Float()
	height = js.Global().Get("innerHeight").Float()

	// Устанавливаем размеры canvas
	canvas.Set("width", width)
	canvas.Set("height", height)

	// Инициализируем частицы
	initParticles()

	// Регистрируем обработчики событий
	canvas.Call("addEventListener", "click", js.FuncOf(handleClick))
	js.Global().Call("addEventListener", "resize", js.FuncOf(handleResize))

	// Регистрируем функции для JavaScript
	registerCallbacks()

	// Выводим информацию в консоль
	js.Global().Get("console").Call("log", "🚀 Go WebAssembly приложение запущено!")
	js.Global().Get("console").Call("log", "📊 Информация о системе:", getSystemInfo())
	js.Global().Get("console").Call("log", "💡 Кликните по canvas для создания взрыва частиц!")
	js.Global().Get("console").Call("log", "🔧 Доступные функции: getSystemInfo(), resetParticles(), addParticles(count)")

	// Запускаем анимацию
	animate()

	// Блокируем завершение программы
	select {}
}
