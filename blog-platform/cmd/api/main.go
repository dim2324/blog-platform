package main

import (
	"blog-platform/internal/handler"
	"blog-platform/internal/middleware"
	"blog-platform/internal/repository"
	"blog-platform/internal/service"
	"blog-platform/pkg/database"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	store := database.NewJSONStore(dataDir)

	userRepo := repository.NewUserRepo(store)
	postRepo := repository.NewPostRepo(store)
	commentRepo := repository.NewCommentRepo(store)

	logChan := make(chan string, 100)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go startLogger(ctx, &wg, logChan)

	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo, logChan)
	commentService := service.NewCommentService(commentRepo, logChan)

	h := handler.NewHandler(userService, postService, commentService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
	mux.Handle("POST /posts", middleware.AuthMiddleware(http.HandlerFunc(h.CreatePost)))
	mux.HandleFunc("GET /posts", h.ListPosts)
	mux.HandleFunc("GET /posts/{id}", h.GetPost)
	mux.Handle("POST /posts/{id}/comments", middleware.AuthMiddleware(http.HandlerFunc(h.CreateComment)))
	mux.HandleFunc("GET /posts/{id}/comments", h.ListComments)
	mux.HandleFunc("GET /health", h.Health)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	cancel()
	wg.Wait()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func startLogger(ctx context.Context, wg *sync.WaitGroup, logChan <-chan string) {
	defer wg.Done()
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Ошибка открытия файла лога: %v", err)
		return
	}
	defer file.Close()

	for {
		select {
		case msg := <-logChan:
			// задержка 1-2 секунды (имитация обработки)
			select {
			case <-time.After(1 * time.Second):
				file.WriteString(msg + "\n")
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
