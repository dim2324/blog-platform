package main

import (
	"log"
	"net/http"

	"tasks-api/internal/handlers"
	httputil "tasks-api/internal/http"
	"tasks-api/internal/storage"
)

func main() {
	// Инициализация in-memory хранилища
	store := storage.NewMemoryStorage()

	// Создание обработчиков
	h := handlers.New(store)

	// Настройка маршрутов с middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", httputil.LoggingMiddleware(h.TasksCollection))
	mux.HandleFunc("/tasks/", httputil.LoggingMiddleware(h.TaskItem))

	// Запуск сервера
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
