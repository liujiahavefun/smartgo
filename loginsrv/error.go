package main

import (
	"errors"
)

var (
	ErrorUserNotFound    error = errors.New("user not found")
	ErrorPasswordInvalid error = errors.New("invalid password")
)
