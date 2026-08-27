package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"blog-platform/internal/model"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.Create(postID, userID, req.Text)
	if err != nil {
		if errors.Is(err, model.ErrEmptyComment) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *Handler) ListComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	comments, err := h.commentService.ListByPostID(postID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
