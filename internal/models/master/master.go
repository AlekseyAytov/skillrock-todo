package master

import (
	"slices"
	"time"

	"github.com/AlekseyAytov/skillrock-todo/internal/models/task"
	"github.com/AlekseyAytov/skillrock-todo/internal/store"
)

// allowedStatuses содержит возможные значения Status
var allowedStatuses = []string{"new", "in_progress", "done"}

// TaskMaster отвечает за работу с объектами Task
type TaskMaster struct {
	store store.ToDoStore
	// можно хранить список задач в мастере
	// taskList []task.Task
}

// NewTaskMaster returns TaskMaster
func NewTaskMaster(st store.ToDoStore) *TaskMaster {
	return &TaskMaster{store: st}
}

func (tm *TaskMaster) createTask(data task.TaskHeads) (*task.Task, error) {
	if err := tm.verifyTitle(data); err != nil {
		return nil, err
	}
	// ID формируется автоматически в tm.store
	t := &task.Task{
		Title:       data.Title,
		Description: data.Description,
		Status:      "new", // при создании новой задачи, ее статус всегда "new"
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return t, nil
}

// func (tm *TaskMaster) generateID() string {
// 	// TODO: изменить формат ID
// 	return fmt.Sprintf("%d", time.Now().UnixNano())
// }

func (tm *TaskMaster) verifyTitle(task task.TaskHeads) error {
	if task.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}

func (tm *TaskMaster) verifyStatus(task task.TaskHeads) error {
	if !slices.Contains(allowedStatuses, task.Status) { // need go 1.21+
		return ErrBadStatus
	}
	return nil
}

// TODO: реализовать при использовании tm.taskList
// TODO: реализовать аргумент функционального типа с возвратом bool
// func (tm *TaskMaster) findBy(id string) (task.Task, error) {
// 	return task.Task{}, nil
// }

// GetAll tasks
func (tm *TaskMaster) GetAll() ([]task.Task, error) {
	return tm.store.GetAll()
}

// Add creates Task from TaskHeads and adds it to store
func (tm *TaskMaster) Add(task task.TaskHeads) error {
	t, err := tm.createTask(task)
	if err != nil {
		return err
	}
	err = tm.store.Add(*t)
	if err != nil {
		return err
	}
	return nil
}

// Delete Task with current ID
func (tm *TaskMaster) Delete(id string) error {
	t, err := tm.store.FindBy(id)
	if err != nil {
		return err
	}
	return tm.store.Delete(t)
}

// UpdateBy takes new data from TaskHeads and updates Task in the store by ID
func (tm *TaskMaster) UpdateBy(id string, newTask task.TaskHeads) error {
	t, err := tm.store.FindBy(id)
	if err != nil {
		return err
	}
	err = tm.verifyStatus(newTask)
	if err != nil {
		return err
	}
	err = tm.verifyTitle(newTask)
	if err != nil {
		return err
	}

	// изменяем найденный объект Task
	t.Title = newTask.Title
	t.Description = newTask.Description
	t.Status = newTask.Status
	t.UpdatedAt = time.Now()

	return tm.store.Update(t)
}
