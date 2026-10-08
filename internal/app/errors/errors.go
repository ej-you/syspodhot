package errors

import "errors"

var (
	ErrNotFound      error = errors.New("record not found") // code 404
	ErrAlreadyExists error = errors.New("already exists")   // code 409
)
