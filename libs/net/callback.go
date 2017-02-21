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
