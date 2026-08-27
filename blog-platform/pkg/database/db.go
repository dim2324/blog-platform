package database

import (
	"blog-platform/internal/model"
	"encoding/json"
	"os"
	"sync"
)

type JSONStore struct {
	mu       sync.RWMutex
	dataDir  string
	users    []model.User
	posts    []model.Post
	comments []model.Comment
}

func NewJSONStore(dataDir string) *JSONStore {
	store := &JSONStore{dataDir: dataDir}
	store.load()
	return store
}

func (s *JSONStore) load() {
	// загрузка users.json, posts.json, comments.json
	// если файл отсутствует - создаём пустой слайс
}

func (s *JSONStore) saveUsers() {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, _ := json.MarshalIndent(s.users, "", "  ")
	os.WriteFile(s.dataDir+"/users.json", data, 0644)
}

// аналогично для posts, comments
