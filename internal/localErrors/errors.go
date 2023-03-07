package localErrors

import "errors"

var (
	ErrAlreadyExists = errors.New("repo: row already exist")
	ErrNotFound      = errors.New("repo: not found order")
	ErrCashNotFound  = errors.New("cache: not found")
)
