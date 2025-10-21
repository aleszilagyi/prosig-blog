package error

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidInput   = errors.New("invalid request input")
)
