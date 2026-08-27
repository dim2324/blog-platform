package repository

import (
	"blog-platform/internal/model"
	"blog-platform/pkg/database"
)

type CommentRepo struct {
	store *database.JSONStore
}

func NewCommentRepo(store *database.JSONStore) *CommentRepo {
	return &CommentRepo{store: store}
}

func (r *CommentRepo) Create(comment *model.Comment) error {
	return r.store.AddComment(comment)
}

func (r *CommentRepo) ListByPostID(postID int) ([]model.Comment, error) {
	comments := r.store.GetComments()
	var result []model.Comment
	for _, c := range comments {
		if c.PostID == postID {
			result = append(result, c)
		}
	}
	return result, nil
}
