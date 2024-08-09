package storage

import (
	"github.com/vadskev/go_final_project/internal/models/task"
)

type TaskRepository interface {
	Create(task *task.Info) (int64, error)
	GetTasks(searchStr string) ([]task.Task, error)
	GetById(id string) (*task.Task, error)
	Update(task *task.Task) error
	Delete(id string) error
}
