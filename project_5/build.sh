#!/bin/bash

echo "🔨 Сборка Go WebAssembly приложения"
echo "==================================="

echo ""
echo "1. Компиляция Go в WebAssembly..."
GOOS=js GOARCH=wasm go build -o main.wasm main.go
if [ $? -ne 0 ]; then
    echo "❌ Ошибка при компиляции в WebAssembly"
    exit 1
fi

echo "✅ WebAssembly файл создан: main.wasm"

echo ""
echo "2. Копирование wasm_exec.js..."
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
if [ $? -ne 0 ]; then
    echo "❌ Ошибка при копировании wasm_exec.js"
    exit 1
fi

echo "✅ wasm_exec.js скопирован"

echo ""
echo "3. Проверка файлов..."
ls -la main.wasm wasm_exec.js index.html

echo ""
echo "✅ Сборка завершена успешно!"
echo ""
echo "📁 Созданные файлы:"
echo "   - main.wasm     - WebAssembly модуль"
echo "   - wasm_exec.js  - JavaScript runtime для Go"
echo "   - index.html    - HTML страница"
echo ""
echo "🚀 Для запуска:"
echo "   1. Запустите HTTP сервер: python -m http.server 8000"
echo "   2. Откройте браузер: http://localhost:8000"
echo ""
echo "💡 Или используйте Go сервер: go run server.go"
