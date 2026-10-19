package errors

import "errors"

var (
	ErrNotFound           = errors.New("record not found")          // code 404
	ErrAlreadyExists      = errors.New("already exists")            // code 409
	ErrInvalidCredentials = errors.New("invalid login or password") // code 400
)
