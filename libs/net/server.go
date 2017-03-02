package net

import (
	"net"
	"time"

	. "smartgo/libs/utils"
)

/*
* liujia: 如果设置了tlsWrapper，就将普通的tcp连接转为TSL连接 (其实就验证key/cert) ，默认啥都不转
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
* server interface for TCP/TLS-TCP server
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

	SetOnCloseCallback(func(Connection))
	GetOnCloseCallback() onCloseFunc

	SetOnErrorCallback(func())
	GetOnErrorCallback() onErrorFunc

	SetOnPacketRecvCallback(callback onPacketRecvFunc)
	GetOnPacketRecvCallback() onPacketRecvFunc

	//下面这俩应该没大用
	SetOnScheduleCallback(time.Duration, func(time.Time, interface{}))
	GetOnScheduleCallback() (time.Duration, onScheduleFunc)

	//SetOnMessageCallback(func(Message, Connection))
	//GetOnMessageCallback() onMessageFunc
}
