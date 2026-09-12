
### Blog Platform API

REST API блог-платформы на Go. Хранение данных — в JSON-файлах, аутентификация — JWT, отложенное логирование действий пользователя.

## Возможности

- 🔐  Регистрация и вход (bcrypt + JWT)
- 📄  Управление постами (создание, список, чтение по ID)
- 💬  Комментарии к постам
- ❤️  Проверка состояния сервиса (`/health`)
- 🕓  Отложенное логирование действий (горутина + канал) в `log.txt` с временными метками
- 🧩  Единый JSON-формат ошибок

## Стек и зависимости

- Go 1.25
- `github.com/golang-jwt/jwt/v5` — JWT
- `golang.org/x/crypto/bcrypt` — хеширование паролей
- Стандартная библиотека `net/http` (Go 1.25 `ServeMux` с метод-роутингом)

## 📁 Структура проекта

blog-platform/

├── cmd/api/main.go
├── data/
├── internal/
│ ├── handler/
│ ├── middleware/
│ ├── model/
│ ├── repository/
│ └── service/
├── pkg/
│ ├── auth/
│ ├── database/
│ └── httpjson/
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md


## ⚙️ Переменные окружения

| Переменная   | Описание                          | По умолчанию |
|--------------|-----------------------------------|--------------|
| `PORT`       | Порт HTTP-сервера                 | `8080`       |
| `DATA_DIR`   | Каталог для JSON-файлов           | `./data`     |
| `JWT_SECRET` | Секрет для подписи JWT            | —            |


🌐 Эндпоинты

🔹 Регистрация

POST /register
Content-Type: application/json

{"email":"user@example.com","username":"user1","password":"secret"}

Ответ 201:
{"id":1,"email":"user@example.com","username":"user1","created_at":"2026-09-12T10:00:00Z"}

Ошибки: 400 (неверный email / пустые поля), 409 (email или username занят).

🔹 Вход

POST /login
Content-Type: application/json

{"email":"user@example.com","password":"secret"}

Ответ 200:

{"token":"eyJhbGciOi..."}

Ошибка: 401 — неверные учётные данные.


🔹 Создание поста (требуется токен)
POST /posts
Authorization: Bearer <TOKEN>
Content-Type: application/json

{"title":"Заголовок","content":"Текст поста"}
Ответ 201 — объект поста. Ошибки: 400 (пустой title/content), 401 (нет токена).

🔹 Список постов
GET /posts
Ответ 200 — массив постов (возможно пустой []).

🔹 Один пост
GET /posts/{id}
Ответ 200 — объект поста, 404 — если не найден.

🔹 Создание комментария (требуется токен)
POST /posts/{id}/comments
Authorization: Bearer <TOKEN>
Content-Type: application/json

{"text":"Отличный пост!"}
Ответ 201 — объект комментария.
Ошибки: 400 (пустой текст), 401, 404 (пост не найден).

🔹 Список комментариев поста
GET /posts/{id}/comments
Ответ 200 — массив комментариев (пустой [], если нет).
Ошибка 404 — если поста нет.

🔹 Health
GET /health
Ответ 200: {"status":"ok"}

❗ Формат ошибок
Все ошибки возвращаются в JSON:

{"error":"описание ошибки"}

🧾 Логирование
Действия пользователя (создание поста/комментария) асинхронно пишутся в log.txt с задержкой 1 секунда и меткой времени:
2026-09-12 10:00:00 user 1 created post 1
2026-09-12 10:01:00 user 1 created comment 1


🧪 Тестирование (curl)

# Регистрация
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","username":"user1","password":"secret"}'

# Вход
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}' | jq -r .token)

# Создание поста
curl -X POST http://localhost:8080/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Привет","content":"Мир"}'

# Комментарий
curl -X POST http://localhost:8080/posts/1/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"text":"Первый!"}'

# Health
curl http://localhost:8080/health



### Сделано с ❤️ на Go
