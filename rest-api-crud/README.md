
## Описание структуры

| Путь | Описание |
|------|----------|
| `go.mod` | Файл модуля Go, определяющий зависимости проекта |
| `main.go` | Главный файл приложения, точка входа |
| `internal/models/models.go` | Определение моделей данных (структуры Task и т.д.) |
| `internal/storage/storage.go` | Интерфейс хранилища данных |
| `internal/storage/memory.go` | Реализация хранилища в памяти (in-memory) |
| `internal/http/middleware.go` | HTTP-мидлвары (логирование, авторизация и т.д.) |
| `handlers/tasks.go` | HTTP-обработчики для работы с задачами (CRUD) |

## Назначение

Проект реализует REST API для управления задачами (CRUD операции) с использованием:
- In-memory хранилище
- Мидлвары для обработки запросов
- Чистая архитектура с разделением на слои

1. Тестирование GET /tasks (получение списка задач)

Корректный сценарий - получение пустого списка:

curl -i -X GET http://localhost:8080/tasks
Ответ:
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:00 GMT
Content-Length: 2

Ошибочный сценарий - неверный метод (POST без тела):

curl -i -X POST http://localhost:8080/tasks
Ответ:
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:01 GMT
Content-Length: 35
{"error":"Invalid JSON format"}

2. Тестирование POST /tasks (создание задачи)
Корректный сценарий - создание задачи:
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Изучить Go","done":false}'

Ответ:
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:02 GMT
Content-Length: 87
{"id":1,"title":"Изучить Go","done":false,"created_at":"2026-07-27T10:00:02.123456Z"}

Ошибочный сценарий - пустой заголовок:
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"","done":false}'
Ответ:
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:03 GMT
Content-Length: 36
{"error":"Title is required"}
