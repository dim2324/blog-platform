package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"blog-platform/internal/service"
	"blog-platform/pkg/httpjson"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail):
			httpjson.WriteError(w, http.StatusBadRequest, "invalid email format")
		case errors.Is(err, service.ErrEmptyFields):
			httpjson.WriteError(w, http.StatusBadRequest, "email, username and password are required")
		case errors.Is(err, service.ErrUserAlreadyExists):
			httpjson.WriteError(w, http.StatusConflict, "email or username already exists")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, user.ToResponse())
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		httpjson.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
