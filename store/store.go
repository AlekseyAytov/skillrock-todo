package store

import (
	"github.com/AlekseyAytov/skillrock-todo/internal/models/task"
)

// ToDoStore some comment
type ToDoStore interface {
	Add(task.Task) error
	FindBy(string) (task.Task, error)
	GetAll() ([]task.Task, error)
	Update(task.Task) error
	Delete(task.Task) error
}
