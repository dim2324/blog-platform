package storage

import (
	"fmt"
	"sync"
	"time"

	"rest-api-crud/internal/models"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	tasks  map[int]models.Task
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *MemoryStorage) List() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *MemoryStorage) Create(task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++
	task.CreatedAt = time.Now()
	s.tasks[task.ID] = task
	return task, nil
}

func (s *MemoryStorage) Get(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}

func (s *MemoryStorage) Update(id int, task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.tasks[id]
	if !exists {
		return models.Task{}, fmt.Errorf("task with id %d not found", id)
	}

	task.ID = id
	task.CreatedAt = existing.CreatedAt
	s.tasks[id] = task
	return task, nil
}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task with id %d not found", id)
	}

	delete(s.tasks, id)
	return nil
}
