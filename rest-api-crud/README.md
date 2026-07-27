
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





# 📋 Тестирование REST API CRUD

## 1. GET /tasks — Получение списка задач

### ✅ Корректный сценарий — получение пустого списка

**Запрос:**

curl -i -X GET http://localhost:8080/tasks
```
**Ответ:**

HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:00 GMT
Content-Length: 2

---

### ❌ Ошибочный сценарий — неверный метод (POST без тела)

**Запрос:**

curl -i -X POST http://localhost:8080/tasks

---
**Ответ:**

HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:01 GMT
Content-Length: 35

{
  "error": "Invalid JSON format"
}

---

## 2. POST /tasks — Создание задачи

### ✅ Корректный сценарий — создание задачи

**Запрос:**

curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Изучить Go","done":false}'
```

**Ответ:**
```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:02 GMT
Content-Length: 87

{
  "id": 1,
  "title": "Изучить Go",
  "done": false,
  "created_at": "2026-07-27T10:00:02.123456Z"
}
```

---

### ❌ Ошибочный сценарий — пустой заголовок

**Запрос:**
```
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"","done":false}'
```

**Ответ:**
```
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:03 GMT
Content-Length: 36

{
  "error": "Title is required"
}
```

---

## 3. GET /tasks/{id} — Получение задачи по ID

### ✅ Корректный сценарий — получение существующей задачи

**Запрос:**
```
curl -i -X GET http://localhost:8080/tasks/1
```

**Ответ:**
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:05 GMT

{
  "id": 1,
  "title": "Изучить Go",
  "done": false,
  "created_at": "2026-07-27T10:00:02.123456Z"
}
```

---

### ❌ Ошибочный сценарий — несуществующий ID

**Запрос:**
```
curl -i -X GET http://localhost:8080/tasks/999
```

**Ответ:**
```
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:06 GMT

{
  "error": "Task not found"
}
```

---

### ❌ Ошибочный сценарий — невалидный ID

**Запрос:**
```
curl -i -X GET http://localhost:8080/tasks/abc
```

**Ответ:**
```
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:07 GMT

{
  "error": "Invalid task ID"
}
```

---

## 4. PUT /tasks/{id} — Обновление задачи

### ✅ Корректный сценарий — обновление существующей задачи

**Запрос:**
```
curl -i -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Изучить Go и REST API","done":true}'
```

**Ответ:**
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:08 GMT

{
  "id": 1,
  "title": "Изучить Go и REST API",
  "done": true,
  "created_at": "2026-07-27T10:00:02.123456Z"
}
```

---

### ❌ Ошибочный сценарий — обновление несуществующей задачи

**Запрос:**
```
curl -i -X PUT http://localhost:8080/tasks/999 \
  -H "Content-Type: application/json" \
  -d '{"title":"Несуществующая задача","done":false}'
```

**Ответ:**
```
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:09 GMT

{
  "error": "Task not found"
}
```

---

### ❌ Ошибочный сценарий — пустой заголовок при обновлении

**Запрос:**
```
curl -i -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"","done":false}'
```

**Ответ:**
```
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:10 GMT

{
  "error": "Title is required"
}
```

---

## 5. DELETE /tasks/{id} — Удаление задачи

### ✅ Корректный сценарий — удаление существующей задачи

**Запрос:**
```bash
curl -i -X DELETE http://localhost:8080/tasks/3
```

**Ответ:**
```
HTTP/1.1 204 No Content
Date: Mon, 27 Jul 2026 10:00:11 GMT
```

> **💡 Примечание:** Статус `204 No Content` означает успешное удаление. Тело ответа отсутствует.

---

### ❌ Ошибочный сценарий — удаление несуществующей задачи

**Запрос:**
```
curl -i -X DELETE http://localhost:8080/tasks/999
```

**Ответ:**
```
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:13 GMT

{
  "error": "Task not found"
}
```

---

## 6. Тестирование неверных HTTP-методов

### ❌ PATCH для коллекции задач

**Запрос:**
```
curl -i -X PATCH http://localhost:8080/tasks
```

**Ответ:**
```
HTTP/1.1 405 Method Not Allowed
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:14 GMT

{
  "error": "Method not allowed"
}
```

---

### ❌ PATCH для конкретной задачи

**Запрос:**
```
curl -i -X PATCH http://localhost:8080/tasks/1
```

**Ответ:**
```
HTTP/1.1 405 Method Not Allowed
Content-Type: application/json
Date: Mon, 27 Jul 2026 10:00:15 GMT

{
  "error": "Method not allowed"
}
```

---

## 📊 Итоговая сводка тестирования

| Метод | Эндпоинт | Сценарий | Статус | Описание |
|-------|----------|----------|:------:|----------|
| `GET` | `/tasks` | Корректный | 🟢 200 | Получен пустой список / список задач |
| `POST` | `/tasks` | Корректный | 🟢 201 | Задача успешно создана |
| `POST` | `/tasks` | Ошибочный | 🟡 400 | Пустой JSON / отсутствует тело запроса |
| `POST` | `/tasks` | Ошибочный | 🟡 400 | Пустой заголовок задачи |
| `GET` | `/tasks/{id}` | Корректный | 🟢 200 | Задача получена по ID |
| `GET` | `/tasks/{id}` | Ошибочный | 🔴 404 | Задача с указанным ID не найдена |
| `GET` | `/tasks/{id}` | Ошибочный | 🟡 400 | Невалидный формат ID |
| `PUT` | `/tasks/{id}` | Корректный | 🟢 200 | Задача успешно обновлена |
| `PUT` | `/tasks/{id}` | Ошибочный | 🔴 404 | Попытка обновить несуществующую задачу |
| `PUT` | `/tasks/{id}` | Ошибочный | 🟡 400 | Пустой заголовок при обновлении |
| `DELETE` | `/tasks/{id}` | Корректный | 🟢 204 | Задача успешно удалена |
| `DELETE` | `/tasks/{id}` | Ошибочный | 🔴 404 | Попытка удалить несуществующую задачу |
| `PATCH` | `/tasks` | Ошибочный | 🟠 405 | Метод не поддерживается для коллекции |
| `PATCH` | `/tasks/{id}` | Ошибочный | 🟠 405 | Метод не поддерживается для элемента |

---
