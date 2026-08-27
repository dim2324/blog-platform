package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidTitle   = errors.New("title cannot be empty")
	ErrInvalidContent = errors.New("content cannot be empty")
	ErrEmptyComment   = errors.New("comment text cannot be empty") // <-- добавьте эту строку
	ErrPostNotFound   = errors.New("post not found")
	ErrUserNotFound   = errors.New("user not found")
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // не сериализуем в JSON
	CreatedAt time.Time `json:"created_at"`
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
