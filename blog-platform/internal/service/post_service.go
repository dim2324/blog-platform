package service

import (
	"errors"
	"fmt"
	"time"

	"blog-platform/internal/model"
	"blog-platform/internal/repository"
)

var (
	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrEmptyContent = errors.New("content cannot be empty")
)

type PostService struct {
	repo    repository.PostRepository
	logChan chan string
}

func NewPostService(repo repository.PostRepository, logChan chan string) *PostService {
	return &PostService{
		repo:    repo,
		logChan: logChan,
	}
}

func (s *PostService) Create(authorID int, title, content string) (*model.Post, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if content == "" {
		return nil, ErrEmptyContent
	}

	post := &model.Post{
		AuthorID:  authorID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(post); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	// Отправляем событие в канал логирования
	s.logChan <- fmt.Sprintf("user %d created post %d", authorID, post.ID)

	return post, nil
}

func (s *PostService) List() ([]model.Post, error) {
	return s.repo.List()
}

func (s *PostService) GetByID(id int) (*model.Post, error) {
	return s.repo.GetByID(id)
}