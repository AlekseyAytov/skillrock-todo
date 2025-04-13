package mock

import (
	"github.com/AlekseyAytov/skillrock-todo/internal/models/task"
	"github.com/AlekseyAytov/skillrock-todo/internal/store"
)

type MockStore struct {
	data []task.Task
}

func NewMockStore() *MockStore {
	m := MockStore{data: []task.Task{}}
	return &m
}

func (m *MockStore) GetAll() ([]task.Task, error) {
	return m.data, nil
}

func (m *MockStore) Add(task task.Task) error {
	m.data = append(m.data, task)
	return nil
}

func (m *MockStore) Delete(task task.Task) error {
	for i, t := range m.data {
		if t.ID == task.ID {
			m.data[i] = m.data[len(m.data)-1]
			m.data = m.data[:len(m.data)-1]
			return nil
		}
	}
	return nil
}

func (m *MockStore) Update(task task.Task) error {
	for _, t := range m.data {
		if t.ID == task.ID {
			t.Title = task.Title
			t.Description = task.Description
			t.Status = task.Status
			t.UpdatedAt = task.UpdatedAt
			return nil
		}
	}
	return store.ErrTaskNotFound
}

func (m *MockStore) FindBy(id string) (task.Task, error) {
	for _, t := range m.data {
		if t.ID == id {
			return t, nil
		}
	}
	return task.Task{}, store.ErrTaskNotFound
}
