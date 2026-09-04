## Структура проекта

```
blog-platform/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── comment_handler.go
│   │   ├── handler.go
│   │   ├── health.go
│   │   └── post_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logging.go
│   ├── model/
│   │   └── models.go
│   ├── repository/
│   │   ├── comment_repo.go
│   │   ├── interfaces.go
│   │   ├── post_repo.go
│   │   └── user_repo.go
│   └── service/
│       ├── comment_service.go
│       ├── post_service.go
│       └── user_service.go
├── pkg/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   └── database/
│       └── db.go
├── data/
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```
