package errs

import "errors"

var (
	ErrInsufficientLimit = errors.New("insufficient limit")
	ErrAccountNotFound   = errors.New("account not found")
)
