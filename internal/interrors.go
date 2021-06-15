package internal

import "errors"

// package internal errors
var (
	ErrInvalidSize = errors.New("invalid size")
	ErrInvalidArg  = errors.New("invalid arg")
	ErrIDNotFound  = errors.New("ID not found")
)
