package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct {
	Store storage.Storage
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(s storage.Storage) *Handler {
	return &Handler{Store: s}
}

func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()
	if tasks == nil {
		tasks = []models.Task{} // Возвращаем пустой массив вместо null
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Валидация
	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	created, err := h.Store.Create(task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request, id int) {
	task, exists := h.Store.Get(id)
	if !exists {
		writeError(w, http.StatusNotFound, "Task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Валидация
	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	updated, err := h.Store.Update(id, task)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.Store.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
