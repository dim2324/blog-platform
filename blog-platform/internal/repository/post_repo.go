package repository

import (
	"blog-platform/internal/model"
	"blog-platform/pkg/database"
)

type PostRepo struct {
	store *database.JSONStore
}

func NewPostRepo(store *database.JSONStore) *PostRepo {
	return &PostRepo{store: store}
}

func (r *PostRepo) Create(post *model.Post) error {
	return r.store.AddPost(post)
}

func (r *PostRepo) GetByID(id int) (*model.Post, error) {
	posts := r.store.GetPosts()
	for i := range posts {
		if posts[i].ID == id {
			return &posts[i], nil
		}
	}
	return nil, model.ErrPostNotFound
}

func (r *PostRepo) List() ([]model.Post, error) {
	return r.store.GetPosts(), nil
}
