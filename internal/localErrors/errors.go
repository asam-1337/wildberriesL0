package localErrors

import "errors"

var (
	ErrAlreadyExists = errors.New("repo: order already exist")
	ErrNotFound      = errors.New("repo: order not found")
)
