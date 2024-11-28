package errors

import (
	"errors"
)

var (
	ErrBadRequest            = errors.New("Bad Request")
	ErrInternalServerError   = errors.New("Internal Server Error")
	ErrInvalidRequest        = errors.New("Invalid request")
	ErrUsernameAlreadyExists = errors.New("Username already exists")
)
