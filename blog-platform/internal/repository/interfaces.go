package repository

import "blog-platform/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	GetByID(id int) (*model.User, error)
}

type PostRepository interface {
	Create(post *model.Post) error
	GetByID(id int) (*model.Post, error)
	List() ([]model.Post, error)
}

type CommentRepository interface {
	Create(comment *model.Comment) error
	ListByPostID(postID int) ([]model.Comment, error)
}
