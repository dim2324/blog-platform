package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidTitle   = errors.New("title cannot be empty")
	ErrInvalidContent = errors.New("content cannot be empty")
	ErrEmptyComment   = errors.New("comment text cannot be empty")
	ErrPostNotFound   = errors.New("post not found")
	ErrUserNotFound   = errors.New("user not found")
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"` // сохраняется в файл
	CreatedAt    time.Time `json:"created_at"`
}

// UserResponse — структура для ответов клиенту (без хеша пароля)
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

type Post struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	AuthorID  int       `json:"author_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
