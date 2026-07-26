# Tasks API

REST API для управления списком задач (to-do list) с хранением данных в памяти.

## Запуск

```bash
# Установка зависимостей
go mod download

# Запуск сервера
go run cmd/server/main.go

Сервер запустится на http://localhost:8080.

API Endpoints
Метод	URL	Описание	Коды ответов
GET	/tasks	Получить список всех задач	200
POST	/tasks	Создать новую задачу	201, 400
GET	/tasks/{id}	Получить задачу по ID	200, 404
PUT	/tasks/{id}	Обновить задачу полностью	200, 400, 404
DELETE	/tasks/{id}	Удалить задачу	204, 404

Примеры использования
1. Создание задачи
bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go", "done": false}'

# Ответ 201:
# {"id":1,"title":"Изучить Go","done":false,"created_at":"2024-01-01T12:00:00Z"}
2. Получение всех задач
bash
curl http://localhost:8080/tasks

# Ответ 200:
# [{"id":1,"title":"Изучить Go","done":false,"created_at":"2024-01-01T12:00:00Z"}]
3. Получение задачи по ID
bash
curl http://localhost:8080/tasks/1

# Ответ 200:
# {"id":1,"title":"Изучить Go","done":false,"created_at":"2024-01-01T12:00:00Z"}
4. Обновление задачи
bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go и gRPC", "done": true}'

# Ответ 200:
# {"id":1,"title":"Изучить Go и gRPC","done":true,"created_at":"2024-01-01T12:00:00Z"}
5. Удаление задачи
bash
curl -X DELETE http://localhost:8080/tasks/1

# Ответ 204 (No Content)

Обработка ошибок
Валидация (400 Bad Request)
bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'

# Ответ 400:
# {"error":"Title is required"}
Не найдено (404 Not Found)
bash
curl http://localhost:8080/tasks/999

# Ответ 404:
# {"error":"Task not found"}
Метод не поддерживается (405)
bash
curl -X PATCH http://localhost:8080/tasks/1

# Ответ 405:
# {"error":"Method not allowed"}
Тестирование
Успешные сценарии
Создание задачи:

bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Тестовая задача", "done": false}'
Обновление задачи:

bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Обновленная задача", "done": true}'
Ошибочные сценарии
Невалидный JSON:

bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d 'invalid json'
Отсутствует title:

bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"done": true}'
Неверный ID:

bash
curl http://localhost:8080/tasks/abc

4)Удаление несуществующей задачи:

bash
curl -X DELETE http://localhost:8080/tasks/999


Особенности реализации
✅ Полностью потокобезопасное in-memory хранилище (sync.RWMutex)

✅ Корректные HTTP статус-коды (200, 201, 204, 400, 404, 405, 500)

✅ Единый формат ошибок в JSON

✅ Валидация входных данных

✅ Логирование всех запросов

✅ Content-Type: application/json для всех ответов

✅ Авто-генерация ID и временных меток



## 8. Тестирование

Запустите сервер и выполните следующие curl команды:

### Успешные сценарии:

```bash
# 1. Создание задачи
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"title":"Задача 1","done":false}'
# Ожидаемый ответ: {"id":1,"title":"Задача 1","done":false,"created_at":"2024-01-01T12:00:00Z"}

# 2. Получение списка
curl http://localhost:8080/tasks
# Ожидаемый ответ: [{"id":1,...}]

# 3. Получение одной задачи
curl http://localhost:8080/tasks/1
# Ожидаемый ответ: {"id":1,"title":"Задача 1",...}

# 4. Обновление задачи
curl -X PUT http://localhost:8080/tasks/1 -H "Content-Type: application/json" -d '{"title":"Обновленная задача","done":true}'
# Ожидаемый ответ: {"id":1,"title":"Обновленная задача","done":true,...}

# 5. Удаление задачи
curl -X DELETE http://localhost:8080/tasks/1
# Ожидаемый ответ: (пустой, статус 204)

Структура проекта
tasks-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   └── tasks.go
│   ├── http/
│   │   └── middleware.go
│   ├── models/
│   │   └── task.go
│   └── storage/
│       ├── storage.go
│       └── memory.go
├── go.mod
├── go.sum
└── README.md