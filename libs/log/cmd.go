package logger

import (
	"fmt"
)

type cmd struct {
	level int
	op    string
	log   string
}

func (c cmd) String() string {
	return fmt.Sprintf("cmd: [%s] --- [%s]", c.op, c.log)
}
