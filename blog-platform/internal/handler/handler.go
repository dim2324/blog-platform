package handler

import (
	"blog-platform/internal/service"
)

type Handler struct {
	userService    *service.UserService
	postService    *service.PostService
	commentService *service.CommentService
}

func NewHandler(us *service.UserService, ps *service.PostService, cs *service.CommentService) *Handler {
	return &Handler{
		userService:    us,
		postService:    ps,
		commentService: cs,
	}
}
