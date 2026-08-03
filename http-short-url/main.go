package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var shortener *URLShortener

func main() {
	shortener = NewURLShortener()

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handleShorten)
	mux.HandleFunc("/", handleRedirect)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		sendError(w, "URL не может быть пустым", http.StatusBadRequest)
		return
	}

	shortID, err := shortener.Shorten(req.URL)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ShortenResponse{
		ShortURL:    shortID,
		OriginalURL: req.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем короткий ID из пути
	shortID := strings.TrimPrefix(r.URL.Path, "/")

	if shortID == "" || shortID == "shorten" {
		sendError(w, "Короткий URL не указан", http.StatusBadRequest)
		return
	}

	originalURL, err := shortener.GetOriginal(shortID)
	if err != nil {
		sendError(w, "URL не найден", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
