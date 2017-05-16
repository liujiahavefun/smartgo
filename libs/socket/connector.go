package socket

import (
	"net"
	"time"

	"smartgo/proto/sessevent"
)

const (
	DEFAULT_CONNECT_RETRY_TIMES = 3
)

type TcpConnector struct {
	*peerBase
	*sessionMgr

	EventDispatcher
	evq EventQueue

	//底层的net.Conn
	conn net.Conn

	//重连间隔时间, 0为不重连
	autoReconnectSec int

	//尝试连接次数
	tryConnTimes int

	//重入锁
	working bool

	//等待关闭的chan
	closeSignal chan bool

	defaultSes Session

	sessionCallbacks *SessionCallback
}

func NewConnector(evq EventQueue) Connector {
	self := &TcpConnector{
		peerBase:    newPeerBase(),
		sessionMgr:  newSessionManager(),
		EventDispatcher: NewEventDispatcher(),
		evq:         evq,
		closeSignal: make(chan bool),
	}

	self.sessionCallbacks = NewSessionCallback(self.onSessionClosedFunc,
		self.onSessionErrorFunc,
		self.onSessionRecvPacketFunc)

	return self
}

//启动，去连接
func (self *TcpConnector) Start(address string) Connector {
	if self.working {
		return self
	}

	go self.connect(address)
	return self
}

//连接，注意重连是会阻塞的，并且连上之后也是阻塞的，所以这个函数要在单独的goroutine里被调用
func (self *TcpConnector) connect(address string) {
	self.working = true

	for {
		self.tryConnTimes++

		//去连接
		conn, err := net.Dial("tcp", address)
		if err != nil {
			ev := newSessionEvent(Event_SessionError, nil, &sessevent.SessionError{Reason:err.Error()})
			self.evq.Post(self, ev)

			if self.tryConnTimes <= DEFAULT_CONNECT_RETRY_TIMES {
				logErrorf("#connect failed(%s) %v", self.name, err.Error())
			}

			if self.tryConnTimes == DEFAULT_CONNECT_RETRY_TIMES {
				logErrorf("(%s) continue reconnecting, but mute log", self.name)
			}

			//没重连就退出
			if self.autoReconnectSec == 0 {
				self.evq.Post(self, newSessionEvent(Event_SessionConnectFailed, nil, &sessevent.SessionConnectFailed{Reason: err.Error()}))
				break
			}

			//有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			//继续连接
			continue
		}

		self.tryConnTimes = 0

		//连上了, 记录连接
		self.conn = conn

		//创建Session
		ses := newClientSession(conn, self, self.sessionCallbacks)
		self.sessionMgr.Add(ses)
		self.defaultSes = ses

		logInfof("#connected(%s) %s sid: %d", self.name, address, ses.id)

		//设置回调
		/*
		ses.onSessionClosedFunc = func(session Session) {
			self.sessionMgr.Remove(session)
			self.closeSignal <- true

			ev := newSessionEvent(Event_SessionClosed, session, &sessevent.SessionClosed{Reason: ""})

			//post断开事件
			self.evq.Post(self, ev)
		}
		ses.onSessionErrorFunc = func(session Session, err error) {

		}
		ses.onSessionRecvPacketFunc = func(session Session, packet *Packet) {
			self.evq.Post(self, &SessionEvent{
				Packet: packet,
				Ses:    session,
			})
		}
		*/

		// 抛出事件
		self.evq.Post(self, NewSessionEvent(Event_SessionConnected, ses, nil))

		//等待连接关闭
		if <-self.closeSignal {
			self.conn = nil

			// 没重连就退出
			if self.autoReconnectSec == 0 {
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue
		}
	}

	self.working = false
}

func (self *TcpConnector) Stop() {
	if self.conn != nil {
		//这个调用会导致session的Close，从而调用了我们设置的OnClose回调，最终self.closeSignal收到信号，关闭
		self.conn.Close()
	}
}

func (self *TcpConnector) Session() Session {
	return self.defaultSes
}

//设置自动重连间隔, 秒为单位，0表示不重连
func (self *TcpConnector) SetAutoReconnectSec(sec int) {
	self.autoReconnectSec = sec
}

func (self *TcpConnector) onSessionClosedFunc(ses Session) {
	self.sessionMgr.Remove(ses)
	self.closeSignal <- true

	//post断开事件
	ev := newSessionEvent(Event_SessionClosed, ses, &sessevent.SessionClosed{})
	self.evq.Post(self, ev)
}

func (self *TcpConnector) onSessionErrorFunc(session Session, err error) {

}

func (self *TcpConnector) onSessionRecvPacketFunc(session Session, packet *Packet) {
	self.evq.Post(self, &SessionEvent{
		Packet: packet,
		Ses:    session,
	})
}