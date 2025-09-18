#!/bin/bash

echo "Сборка кроссплатформенного CLI-приложения"
echo "=========================================="

echo ""
echo "Сборка для Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o finder_windows.exe main.go
if [ $? -ne 0 ]; then
    echo "Ошибка при сборке для Windows"
    exit 1
fi

echo ""
echo "Сборка для Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o finder_linux main.go
if [ $? -ne 0 ]; then
    echo "Ошибка при сборке для Linux"
    exit 1
fi

echo ""
echo "Сборка для macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o finder_mac_arm64 main.go
if [ $? -ne 0 ]; then
    echo "Ошибка при сборке для macOS ARM64"
    exit 1
fi

echo ""
echo "Сборка для macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o finder_mac_amd64 main.go
if [ $? -ne 0 ]; then
    echo "Ошибка при сборке для macOS AMD64"
    exit 1
fi

echo ""
echo "✅ Сборка завершена успешно!"
echo ""
echo "Созданные файлы:"
ls -la finder_*
echo ""
echo "Размеры файлов:"
for file in finder_*; do
    if [ -f "$file" ]; then
        size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
        echo "$file: $size байт"
    fi
done
echo ""
echo "Теперь вы можете запускать эти файлы на соответствующих платформах!"
