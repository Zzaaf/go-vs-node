#!/bin/bash

echo "🧪 Тестирование REST API"
echo "========================"

BASE_URL="http://localhost:8080"

echo ""
echo "1. Получение информации об API..."
curl -s "$BASE_URL/" | jq '.'

echo ""
echo "2. Получение всех задач..."
curl -s "$BASE_URL/tasks" | jq '.'

echo ""
echo "3. Создание новой задачи..."
curl -s -X POST "$BASE_URL/create" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Тестовая задача",
    "description": "Задача создана через API",
    "priority": "high"
  }' | jq '.'

echo ""
echo "4. Получение задачи по ID (ID=1)..."
curl -s "$BASE_URL/tasks/1" | jq '.'

echo ""
echo "5. Обновление задачи (ID=1)..."
curl -s -X PUT "$BASE_URL/tasks/1" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "priority": "low"
  }' | jq '.'

echo ""
echo "6. Получение всех задач после обновления..."
curl -s "$BASE_URL/tasks" | jq '.'

echo ""
echo "7. Удаление задачи (ID=1)..."
curl -s -X DELETE "$BASE_URL/tasks/1" | jq '.'

echo ""
echo "8. Получение всех задач после удаления..."
curl -s "$BASE_URL/tasks" | jq '.'

echo ""
echo "✅ Тестирование завершено!"
