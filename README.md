📚 Примеры использования
1. Создание задачи
bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go", "done": false}'
Ответ (201 Created):

json
{
  "id": 1,
  "title": "Изучить Go",
  "done": false,
  "created_at": "2024-01-01T12:00:00Z"
}
2. Получение всех задач
bash
curl http://localhost:8080/tasks
Ответ (200 OK):

json
[
  {
    "id": 1,
    "title": "Изучить Go",
    "done": false,
    "created_at": "2024-01-01T12:00:00Z"
  }
]
3. Получение задачи по ID
bash
curl http://localhost:8080/tasks/1
Ответ (200 OK):

json
{
  "id": 1,
  "title": "Изучить Go",
  "done": false,
  "created_at": "2024-01-01T12:00:00Z"
}
4. Обновление задачи
bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go и gRPC", "done": true}'
Ответ (200 OK):

json
{
  "id": 1,
  "title": "Изучить Go и gRPC",
  "done": true,
  "created_at": "2024-01-01T12:00:00Z"
}
5. Удаление задачи
bash
curl -X DELETE http://localhost:8080/tasks/1
Ответ: 204 No Content (пустое тело)

❌ Обработка ошибок
Все ошибки возвращаются в формате JSON с полем error.

Валидация (400 Bad Request)
bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'
Ответ:

json
{
  "error": "Title is required"
}
Не найдено (404 Not Found)
bash
curl http://localhost:8080/tasks/999
Ответ:

json
{
  "error": "Task not found"
}
Метод не поддерживается (405 Method Not Allowed)
bash
curl -X PATCH http://localhost:8080/tasks/1
Ответ:

json
{
  "error": "Method not allowed"
}
🧪 Тестирование
Ниже приведены готовые сценарии для ручного тестирования API.

✅ Успешные сценарии
bash
# 1. Создание задачи
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Тестовая задача", "done": false}'

# 2. Обновление задачи
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Обновленная задача", "done": true}'
❌ Ошибочные сценарии
bash
# Невалидный JSON
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d 'invalid json'

# Отсутствует поле title
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"done": true}'

# Неверный формат ID
curl http://localhost:8080/tasks/abc

# Удаление несуществующей задачи
curl -X DELETE http://localhost:8080/tasks/999
⚙️ Особенности реализации
✅ Полностью потокобезопасное in-memory хранилище (sync.RWMutex)

✅ Корректные HTTP статус-коды (200, 201, 204, 400, 404, 405, 500)

✅ Единый формат ошибок в JSON

✅ Валидация входных данных

✅ Логирование всех запросов

✅ Content-Type: application/json для всех ответов

✅ Авто-генерация ID и временных меток

📁 Структура проекта
text
tasks-api/
├── cmd/
│   └── server/
│       └── main.go            # Точка входа в приложение
├── internal/
│   ├── handlers/
│   │   └── tasks.go           # HTTP-обработчики запросов
│   ├── http/
│   │   └── middleware.go      # Middleware (логирование и т.д.)
│   ├── models/
│   │   └── task.go            # Модель данных задачи
│   └── storage/
│       ├── storage.go         # Интерфейс хранилища
│       └── memory.go          # Реализация in-memory хранилища
├── go.mod                     # Файл модуля Go
├── go.sum                     # Контрольные суммы зависимостей
└── README.md                  # Документация проекта
