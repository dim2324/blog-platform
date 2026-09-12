package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"blog-platform/internal/model"
	"blog-platform/pkg/httpjson"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		httpjson.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.postService.Create(userID, req.Title, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidTitle):
			httpjson.WriteError(w, http.StatusBadRequest, "title cannot be empty")
		case errors.Is(err, model.ErrInvalidContent):
			httpjson.WriteError(w, http.StatusBadRequest, "content cannot be empty")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, post)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.List()
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if posts == nil {
		posts = []model.Post{}
	}
	httpjson.WriteJSON(w, http.StatusOK, posts)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.postService.GetByID(id)
	if err != nil {
		if errors.Is(err, model.ErrPostNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "post not found")
		} else {
			httpjson.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, post)
}
