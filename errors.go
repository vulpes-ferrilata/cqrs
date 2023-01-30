package cqrs

import "errors"

var (
	ErrHandlerNotFound error = errors.New("handler not found")
)
