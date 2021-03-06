package utils

import (
	"errors"
)

var (
	ErrorNilKey      		error = errors.New("Nil key")
	ErrorNilValue    		error = errors.New("Nil value")
	ErrorKeyAlreadyExist 	error = errors.New("Key already exist")
	ErrorNotHashable 		error = errors.New("Not hashable")
	ErrorWouldBlock  		error = errors.New("Would block")
)
