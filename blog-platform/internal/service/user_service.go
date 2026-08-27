package service

import (
	"errors"
	"time"

	"blog-platform/internal/model"
	"blog-platform/internal/repository"
	"blog-platform/pkg/auth"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, username, password string) (*model.User, error) {
	// Проверка уникальности email и username
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, ErrUserAlreadyExists
	}
	if _, err := s.repo.FindByUsername(username); err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Хэширование пароля
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
