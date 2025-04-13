package store

import "errors"

// ErrTaskNotFound indicates absence of Task
var ErrTaskNotFound = errors.New("task not found")
