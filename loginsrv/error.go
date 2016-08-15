package main

import (
	"errors"
)

var (
	ErrorUserNotFound    error = errors.New("user not found")
	ErrorInvalidPassword error = errors.New("invalid password")
	ErrorInvalidToken    error = errors.New("invalid token")
)
