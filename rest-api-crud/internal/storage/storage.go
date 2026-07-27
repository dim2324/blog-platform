package storage

import "rest-api-crud/internal/models"

type Storage interface {
	List() []models.Task
	Create(task models.Task) (models.Task, error)
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}
