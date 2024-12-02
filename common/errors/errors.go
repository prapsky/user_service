package errors

import (
	"errors"
)

var (
	ErrBadRequest            = errors.New("Bad Request")
	ErrInternalServerError   = errors.New("Internal Server Error")
	ErrInvalidRequest        = errors.New("Invalid request")
	ErrUsernameAlreadyExists = errors.New("Username already exists")
	ErrAccountNotRegistered  = errors.New("Account not registered")
	ErrInvalidToken          = errors.New("User token not valid")
	ErrUnexpectedSigning     = errors.New("Unexpected signing")
	ErrIncorrectPassword     = errors.New("Incorrect Phone Number or Password")
	ErrUnauthorizedUser      = errors.New("Unauthorized user")
)
