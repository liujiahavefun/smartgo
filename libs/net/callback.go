package net

import (
	"smartgo/libs/pool"
	"time"
)

type onConnectFunc func(Connection) bool
type onMessageFunc func(Message, Connection)
type onCloseFunc func(Connection)
type onErrorFunc func()
type onScheduleFunc func(time.Time, interface{})

type HandlerProc func()
type onPacketRecvFunc func(Connection, *pool.Buffer) (HandlerProc, bool)

var DBInitializer = func() (interface{}, error) {
	return nil, ErrorNotImplemented
}

func SetDBInitializer(fn func() (interface{}, error)) {
	DBInitializer = fn
}
