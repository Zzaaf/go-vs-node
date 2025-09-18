/**
 * Простой HTTP-сервер на Node.js для демонстрации производительности
 * Этот сервер использует Event Loop и может блокироваться при долгих операциях
 * 
 * @author Go Master Class
 * @version 2.0
 */

const http = require('http');

// ================================
// КОНСТАНТЫ И КОНФИГУРАЦИЯ
// ================================
const CONFIG = {
    PORT: 3000,
    SLOW_OPERATION_DURATION: 10000, // 10 секунд
    ROUTES: {
        ROOT: '/',
        SLOW: '/slow'
    },
    CORS_HEADERS: {
        'Access-Control-Allow-Origin': '*',
        'Content-Type': 'application/json'
    }
};

// ================================
// УТИЛИТЫ И ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ================================

/**
 * ANSI цветовые коды и символы для форматирования консоли
 */
const COLORS = {
    RESET: '\x1b[0m',
    BRIGHT: '\x1b[1m',
    DIM: '\x1b[2m',
    RED: '\x1b[31m',
    GREEN: '\x1b[32m',
    YELLOW: '\x1b[33m',
    BLUE: '\x1b[34m',
    MAGENTA: '\x1b[35m',
    CYAN: '\x1b[36m',
    WHITE: '\x1b[37m',
    BG_RED: '\x1b[41m',
    BG_GREEN: '\x1b[42m',
    BG_YELLOW: '\x1b[43m',
    BG_BLUE: '\x1b[44m'
};

const SYMBOLS = {
    SERVER: '[SERVER]',
    ERROR: '[ERROR]',
    SUCCESS: '[SUCCESS]',
    WARNING: '[WARNING]',
    INFO: '[INFO]',
    CLOCK: '[CLOCK]',
    LIGHTNING: '[FAST]',
    TURTLE: '[SLOW]',
    ROCKET: '[ROCKET]',
    GEAR: '[GEAR]',
    CHECK: '[OK]',
    CROSS: '[FAIL]',
    ARROW: '->',
    BULLET: '•',
    DASH: '─',
    PIPE: '│'
};

/**
 * Форматирует текущее время для логов
 * @returns {string} Отформатированное время
 */
function getCurrentTime() {
    return new Date().toISOString();
}

/**
 * Логирует сообщение с цветным форматированием и временной меткой
 * @param {string} type - Тип сообщения (ERROR, SUCCESS, INFO, etc.)
 * @param {string} message - Текст сообщения
 */
function logWithFormat(type, message) {
    const timestamp = getCurrentTime();
    const symbol = SYMBOLS[type] || SYMBOLS.INFO;

    let color = COLORS.WHITE;
    switch (type) {
        case 'ERROR':
            color = COLORS.RED;
            break;
        case 'SUCCESS':
            color = COLORS.GREEN;
            break;
        case 'WARNING':
            color = COLORS.YELLOW;
            break;
        case 'INFO':
            color = COLORS.CYAN;
            break;
        case 'SERVER':
            color = COLORS.MAGENTA;
            break;
        case 'CLOCK':
        case 'TURTLE':
            color = COLORS.BLUE;
            break;
        case 'LIGHTNING':
            color = COLORS.GREEN;
            break;
        default:
            color = COLORS.WHITE;
    }

    console.log(`${color}${symbol}${COLORS.RESET} ${COLORS.DIM}[${timestamp}]${COLORS.RESET} ${message}`);
}

/**
 * Имитирует долгую операцию (например, запрос к базе данных)
 * @returns {string} Результат операции
 */
function simulateLongOperation() {
    logWithFormat('CLOCK', 'Начало обработки медленного запроса');

    const start = Date.now();
    while (Date.now() - start < CONFIG.SLOW_OPERATION_DURATION) {
        // Имитируем долгую операцию
    }

    logWithFormat('TURTLE', 'Завершение обработки медленного запроса');
    return "Долгая операция завершена!";
}

/**
 * Устанавливает CORS заголовки для ответа
 * @param {http.ServerResponse} res - Объект ответа сервера
 */
function setCorsHeaders(res) {
    Object.entries(CONFIG.CORS_HEADERS).forEach(([key, value]) => {
        res.setHeader(key, value);
    });
}

/**
 * Отправляет JSON ответ клиенту
 * @param {http.ServerResponse} res - Объект ответа сервера
 * @param {number} statusCode - HTTP статус код
 * @param {Object} data - Данные для отправки
 */
function sendJsonResponse(res, statusCode, data) {
    res.writeHead(statusCode);
    res.end(JSON.stringify(data, null, 2));
}

// ================================
// ОБРАБОТЧИКИ МАРШРУТОВ
// ================================

/**
 * Обрабатывает корневой маршрут
 * @param {http.ServerResponse} res - Объект ответа сервера
 */
function handleRootRoute(res) {
    const responseData = {
        message: "Привет от Node.js сервера!",
        timestamp: getCurrentTime()
    };

    sendJsonResponse(res, 200, responseData);
}

/**
 * Обрабатывает медленный маршрут
 * @param {http.ServerResponse} res - Объект ответа сервера
 */
function handleSlowRoute(res) {
    const result = simulateLongOperation();

    const responseData = {
        message: result,
        timestamp: getCurrentTime(),
        note: "Этот запрос заблокировал Event Loop на 10 секунд!"
    };

    sendJsonResponse(res, 200, responseData);
}

/**
 * Обрабатывает несуществующие маршруты
 * @param {http.ServerResponse} res - Объект ответа сервера
 */
function handleNotFoundRoute(res) {
    const responseData = {
        message: "Страница не найдена",
        timestamp: getCurrentTime()
    };

    sendJsonResponse(res, 404, responseData);
}

// ================================
// ОСНОВНОЙ ОБРАБОТЧИК ЗАПРОСОВ
// ================================

/**
 * Создает и настраивает HTTP-сервер
 * @returns {http.Server} Настроенный сервер
 */
function createServer() {
    return http.createServer((req, res) => {
        // Устанавливаем CORS заголовки
        setCorsHeaders(res);

        // Логируем входящий запрос
        logWithFormat('INFO', `Входящий запрос: ${req.method} ${req.url}`);

        // Маршрутизация
        try {
            switch (req.url) {
                case CONFIG.ROUTES.ROOT:
                    handleRootRoute(res);
                    break;
                case CONFIG.ROUTES.SLOW:
                    handleSlowRoute(res);
                    break;
                default:
                    handleNotFoundRoute(res);
                    break;
            }
        } catch (error) {
            logWithFormat('ERROR', `Ошибка при обработке запроса: ${error.message}`);
            sendJsonResponse(res, 500, {
                error: "Внутренняя ошибка сервера",
                timestamp: getCurrentTime(),
                status: "error"
            });
        }
    });
}

// ================================
// ЗАПУСК СЕРВЕРА
// ================================

// Функция для запуска сервера
function startServer() {
    // Создаем сервер
    const server = createServer();

    // Запускаем сервер на указанном порту
    server.listen(CONFIG.PORT, (err) => {
        if (err) {
            logWithFormat('ERROR', `Ошибка при запуске сервера: ${err.message}`);
            process.exit(1);
        }

        // Выводим информацию о запуске сервера
        console.log('\n' + COLORS.CYAN + SYMBOLS.DASH.repeat(60) + COLORS.RESET);
        logWithFormat('SERVER', `Node.js сервер запущен на http://localhost:${CONFIG.PORT}`);
        console.log(COLORS.CYAN + SYMBOLS.DASH.repeat(60) + COLORS.RESET);

        console.log(`\n${COLORS.YELLOW}${SYMBOLS.BULLET} Доступные маршруты:${COLORS.RESET}`);
        console.log(`   ${COLORS.GREEN}${SYMBOLS.LIGHTNING}${COLORS.RESET} GET / - быстрый ответ`);
        console.log(`   ${COLORS.BLUE}${SYMBOLS.TURTLE}${COLORS.RESET} GET /slow - медленный ответ (10 сек)`);

        console.log(`\n${COLORS.YELLOW}${SYMBOLS.WARNING}${COLORS.RESET} Проблема: Event Loop блокируется при долгих операциях!`);
        console.log(COLORS.CYAN + SYMBOLS.DASH.repeat(60) + COLORS.RESET + '\n');
    });

    // Обработка ошибок сервера
    server.on('error', (err) => {
        logWithFormat('ERROR', `Ошибка сервера: ${err.message}`);
    });

    // Обработка сигналов завершения работы сервера
    process.on('SIGTERM', () => {
        logWithFormat('INFO', 'Получен сигнал SIGTERM, завершение работы...');
        server.close(() => {
            logWithFormat('SUCCESS', 'Сервер успешно остановлен');
            process.exit(0);
        });
    });

    // Обработка сигналов завершения работы сервера
    process.on('SIGINT', () => {
        logWithFormat('INFO', 'Получен сигнал SIGINT, завершение работы...');
        server.close(() => {
            logWithFormat('SUCCESS', 'Сервер успешно остановлен');
            process.exit(0);
        });
    });
}

// Запускаем сервер
startServer();
