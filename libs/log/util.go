package logger

import (
	"fmt"
	"runtime"
)

func callstack() (msg string) {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		f := runtime.FuncForPC(pc)
		msg += fmt.Sprintf("frame = %v, file = %v, line = %v, func = %v\n", skip, file, line, f.Name())
	}
	return msg
}
