@echo off
echo 🔨 Сборка Go WebAssembly приложения
echo ===================================

echo.
echo 1. Компиляция Go в WebAssembly...
set GOOS=js
set GOARCH=wasm
go build -o main.wasm main.go
if %errorlevel% neq 0 (
    echo ❌ Ошибка при компиляции в WebAssembly
    exit /b 1
)

echo ✅ WebAssembly файл создан: main.wasm

echo.
echo 2. Копирование wasm_exec.js...
for /f "tokens=*" %%i in ('go env GOROOT') do set GOROOT=%%i
copy "%GOROOT%\misc\wasm\wasm_exec.js" .
if %errorlevel% neq 0 (
    echo ❌ Ошибка при копировании wasm_exec.js
    exit /b 1
)

echo ✅ wasm_exec.js скопирован

echo.
echo 3. Проверка файлов...
dir main.wasm wasm_exec.js index.html

echo.
echo ✅ Сборка завершена успешно!
echo.
echo 📁 Созданные файлы:
echo    - main.wasm     - WebAssembly модуль
echo    - wasm_exec.js  - JavaScript runtime для Go
echo    - index.html    - HTML страница
echo.
echo 🚀 Для запуска:
echo    1. Запустите HTTP сервер: python -m http.server 8000
echo    2. Откройте браузер: http://localhost:8000
echo.
echo 💡 Или используйте Go сервер: go run server.go
