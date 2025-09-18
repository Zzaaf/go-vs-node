@echo off
echo Сборка кроссплатформенного CLI-приложения
echo ==========================================

echo.
echo Сборка для Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -o finder_windows.exe main.go
if %errorlevel% neq 0 (
    echo Ошибка при сборке для Windows
    exit /b 1
)

echo.
echo Сборка для Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -o finder_linux main.go
if %errorlevel% neq 0 (
    echo Ошибка при сборке для Linux
    exit /b 1
)

echo.
echo Сборка для macOS (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -o finder_mac_arm64 main.go
if %errorlevel% neq 0 (
    echo Ошибка при сборке для macOS ARM64
    exit /b 1
)

echo.
echo Сборка для macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -o finder_mac_amd64 main.go
if %errorlevel% neq 0 (
    echo Ошибка при сборке для macOS AMD64
    exit /b 1
)

echo.
echo ✅ Сборка завершена успешно!
echo.
echo Созданные файлы:
dir finder_*.*
echo.
echo Размеры файлов:
for %%f in (finder_*.*) do (
    echo %%~nf: %%~zf байт
)
echo.
echo Теперь вы можете запускать эти файлы на соответствующих платформах!
