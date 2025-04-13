package master

import "errors"

// ErrBadStatus indicates bad status
var ErrBadStatus = errors.New("bad status")

// ErrEmptyTitle indicates bad title
var ErrEmptyTitle = errors.New("title is empty")
