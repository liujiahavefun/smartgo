package utils

import (
	"errors"
)

var (
	ErrorNilKey      error = errors.New("Nil key")
	ErrorNilValue    error = errors.New("Nil value")
	ErrorNotHashable error = errors.New("Not hashable")
	ErrorWouldBlock  error = errors.New("Would block")
)
