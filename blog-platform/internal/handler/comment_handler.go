package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"blog-platform/internal/model"
	"blog-platform/pkg/httpjson"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		httpjson.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comment, err := h.commentService.Create(postID, userID, req.Text)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyComment):
			httpjson.WriteError(w, http.StatusBadRequest, "comment text cannot be empty")
		case errors.Is(err, model.ErrPostNotFound):
			httpjson.WriteError(w, http.StatusNotFound, "post not found")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, comment)
}

func (h *Handler) ListComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	// проверяем существование поста
	if _, err := h.postService.GetByID(postID); err != nil {
		if errors.Is(err, model.ErrPostNotFound) {
			httpjson.WriteError(w, http.StatusNotFound, "post not found")
			return
		}
		httpjson.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	comments, err := h.commentService.ListByPostID(postID)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if comments == nil {
		comments = []model.Comment{}
	}
	httpjson.WriteJSON(w, http.StatusOK, comments)
}
