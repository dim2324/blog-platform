package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"blog-platform/internal/model"
)

type JSONStore struct {
	mu       sync.RWMutex
	dataDir  string
	users    []model.User
	posts    []model.Post
	comments []model.Comment
}

func NewJSONStore(dataDir string) *JSONStore {
	store := &JSONStore{
		dataDir:  dataDir,
		users:    []model.User{},
		posts:    []model.Post{},
		comments: []model.Comment{},
	}
	store.load()
	return store
}

func (s *JSONStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// users
	usersFile := filepath.Join(s.dataDir, "users.json")
	if data, err := os.ReadFile(usersFile); err == nil {
		json.Unmarshal(data, &s.users)
	}

	// posts
	postsFile := filepath.Join(s.dataDir, "posts.json")
	if data, err := os.ReadFile(postsFile); err == nil {
		json.Unmarshal(data, &s.posts)
	}

	// comments
	commentsFile := filepath.Join(s.dataDir, "comments.json")
	if data, err := os.ReadFile(commentsFile); err == nil {
		json.Unmarshal(data, &s.comments)
	}
}

func (s *JSONStore) saveUsers() {
	data, _ := json.MarshalIndent(s.users, "", "  ")
	os.WriteFile(filepath.Join(s.dataDir, "users.json"), data, 0644)
}

func (s *JSONStore) savePosts() {
	data, _ := json.MarshalIndent(s.posts, "", "  ")
	os.WriteFile(filepath.Join(s.dataDir, "posts.json"), data, 0644)
}

func (s *JSONStore) saveComments() {
	data, _ := json.MarshalIndent(s.comments, "", "  ")
	os.WriteFile(filepath.Join(s.dataDir, "comments.json"), data, 0644)
}

// Методы доступа для пользователей
func (s *JSONStore) GetUsers() []model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.User(nil), s.users...)
}

func (s *JSONStore) AddUser(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = len(s.users) + 1
	s.users = append(s.users, *user)
	s.saveUsers()
	return nil
}

// Методы для постов
func (s *JSONStore) GetPosts() []model.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Post(nil), s.posts...)
}

func (s *JSONStore) AddPost(post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	post.ID = len(s.posts) + 1
	s.posts = append(s.posts, *post)
	s.savePosts()
	return nil
}

// Методы для комментариев
func (s *JSONStore) GetComments() []model.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Comment(nil), s.comments...)
}

func (s *JSONStore) AddComment(comment *model.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	comment.ID = len(s.comments) + 1
	s.comments = append(s.comments, *comment)
	s.saveComments()
	return nil
}
