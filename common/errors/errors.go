package errors

import (
	"errors"
)

var (
	ErrUsernameAlreadyExists = errors.New("Username already exists")
)
