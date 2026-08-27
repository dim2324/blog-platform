package service

import (
	"fmt"
	"time"

	"blog-platform/internal/model"
	"blog-platform/internal/repository"
)

type CommentService struct {
	repo    repository.CommentRepository
	logChan chan string
}

func NewCommentService(repo repository.CommentRepository, logChan chan string) *CommentService {
	return &CommentService{repo: repo, logChan: logChan}
}

func (s *CommentService) Create(postID, authorID int, text string) (*model.Comment, error) {
	if text == "" {
		return nil, model.ErrEmptyComment
	}

	comment := &model.Comment{
		PostID:    postID,
		AuthorID:  authorID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	s.logChan <- fmt.Sprintf("user %d created comment %d", authorID, comment.ID)
	return comment, nil
}

func (s *CommentService) ListByPostID(postID int) ([]model.Comment, error) {
	return s.repo.ListByPostID(postID)
}
