package net

import (
	"net"
	"time"

	. "smartgo/libs/utils"
)

/*
* liujia: return tls conn if tlsWrapper is configured, return origin conn otherwise
 */
func init() {
	tlsWrapper = func(conn net.Conn) net.Conn {
		return conn
	}
}

var (
	tlsWrapper func(net.Conn) net.Conn
)

func setTLSWrapper(wrapper func(conn net.Conn) net.Conn) {
	tlsWrapper = wrapper
}

/*
* server interface for TCP/TLS_TCP server
 */
type Server interface {
	Start()
	Close()
	IsRunning() bool

	GetServerAddress() string
	GetAllConnections() *ConcurrentMap
	GetTimingWheel() *TimingWheel
	GetWorkerPool() *WorkerPool

	SetOnConnectCallback(func(Connection) bool)
	GetOnConnectCallback() onConnectFunc

	//SetOnMessageCallback(func(Message, Connection))
	//GetOnMessageCallback() onMessageFunc

	SetOnCloseCallback(func(Connection))
	GetOnCloseCallback() onCloseFunc

	SetOnErrorCallback(func())
	GetOnErrorCallback() onErrorFunc

	SetOnScheduleCallback(time.Duration, func(time.Time, interface{}))
	GetOnScheduleCallback() (time.Duration, onScheduleFunc)

	SetOnPacketRecvCallback(callback onPacketRecvFunc)
	GetOnPacketRecvCallback() onPacketRecvFunc
}
