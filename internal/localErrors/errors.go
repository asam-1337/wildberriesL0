package localErrors

import "errors"

var (
	ErrAlreadyExists = errors.New("row already exist")
	ErrNotFound      = errors.New("not found order")
)
