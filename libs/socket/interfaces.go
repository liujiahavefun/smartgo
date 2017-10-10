package socket

/*
* 本文件定义常用的接口
*/

type Peer interface {
	//名字
	SetName(string)
	Name() string

	//Session最大包大小, 超过这个数字, 接收视为错误, 断开连接
	SetMaxPacketSize(size int)
	MaxPacketSize() int
}

type Server interface {
	Peer

	//事件
	EventDispatcher

	//连接管理
	SessionManager

	//启动服务，侦听接口
	Start(address string) Server

	//停止
	Stop()

	//是否在运行
	IsRunning() bool

	//监听的地址like tcp:127.0.0.1:9100
	GetAddress() string
}

type Connector interface {
	Peer

	//事件
	EventDispatcher

	//启动，连接server
	Start(address string) Connector

	//停止
	Stop()

	//连接后的Session
	Session() Session

	//自动重连间隔, 0表示不重连, 默认不重连
	SetAutoReconnectSec(sec int)
}

type Session interface {
	//发包
	Send(interface{})

	//直接发送封包
	RawSend(*Packet)

	//断开
	Close()

	//标示ID
	SetID(int64)
	GetID() int64

	//归属端
	FromPeer() Peer
	SelfPeer() Peer

	//自定义属性
	GetParam(string) interface{}
	SetParam(string, interface{})
}

type SessionManager interface {
	//获取一个连接
	GetSession(int64) Session

	//遍历连接
	VisitSession(func(Session) bool)

	//连接数量
	SessionCount() int
}

type ServerSession interface {
	Session
}

type ClientSession interface {
	Session
}