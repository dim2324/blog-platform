package repository

import (
	"errors"

	"blog-platform/internal/model"
	"blog-platform/pkg/database"
)

type UserRepo struct {
	store *database.JSONStore
}

func NewUserRepo(store *database.JSONStore) *UserRepo {
	return &UserRepo{store: store}
}

func (r *UserRepo) Create(user *model.User) error {
	return r.store.AddUser(user)
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	users := r.store.GetUsers()
	for i := range users {
		if users[i].Email == email {
			return &users[i], nil
		}
	}
	return nil, model.ErrUserNotFound
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	users := r.store.GetUsers()
	for i := range users {
		if users[i].Username == username {
			return &users[i], nil
		}
	}
	return nil, model.ErrUserNotFound
}

func (r *UserRepo) GetByID(id int) (*model.User, error) {
	users := r.store.GetUsers()
	for i := range users {
		if users[i].ID == id {
			return &users[i], nil
		}
	}
	return nil, model.ErrUserNotFound
}

// для избежания ошибки "declared and not used: errors"
var _ = errors.New
