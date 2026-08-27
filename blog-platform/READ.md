blog-platform/
├── cmd/api/
│   └── main.go
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── post_handler.go
│   │   ├── comment_handler.go
│   │   └── health.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logging.go
│   ├── model/
│   │   └── models.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── user_repo.go
│   │   ├── post_repo.go
│   │   └── comment_repo.go
│   └── service/
│       ├── user_service.go
│       ├── post_service.go
│       └── comment_service.go
├── pkg/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   └── database/
│       └── db.go
├── data/               # JSON-файлы создаются автоматически
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md