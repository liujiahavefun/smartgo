package net

import (
	"net"
	"sync"
	"time"

	"github.com/golang/glog"

	"smartgo/libs/pool"
	. "smartgo/libs/utils"
)

// implements Server Side Connection
type ServerConnection struct {
	netid         int64
	name          string
	heartBeat     int64
	extraData     interface{}
	owner         *TCPServer
	isClosed      *AtomicBoolean
	once          sync.Once
	pendingTimers []int64
	conn          net.Conn
	messageCodec  Codec
	finish        *sync.WaitGroup

	messageSendChan chan []byte
	handlerRecvChan chan MessageHandler
	packetRecvChan  chan *pool.Buffer
	closeConnChan   chan struct{}
	timeOutChan     chan *OnTimeOut

	onConnect onConnectFunc
	onMessage onMessageFunc
	onClose   onCloseFunc
	onError   onErrorFunc
	onPacket  onPacketRecvFunc
}

func NewServerConnection(netid int64, server *TCPServer, conn net.Conn) Connection {
	serverConn := &ServerConnection{
		netid:           netid,
		name:            conn.RemoteAddr().String(),
		heartBeat:       time.Now().UnixNano(),
		owner:           server,
		isClosed:        NewAtomicBoolean(false),
		pendingTimers:   []int64{},
		conn:            conn,
		messageCodec:    NewFixedLengthHeaderCodec(HEADER_SIZE, MIN_PACKET_SIZE, MAX_PACKET_SIZE, GLOBAL_BINARY_BYTE_ORDER),
		finish:          &sync.WaitGroup{},
		messageSendChan: make(chan []byte, 1024),
		handlerRecvChan: make(chan MessageHandler, 1024),
		packetRecvChan:  make(chan *pool.Buffer, 1024),
		closeConnChan:   make(chan struct{}),
		timeOutChan:     make(chan *OnTimeOut),
	}

	serverConn.SetOnConnectCallback(server.onConnect)
	serverConn.SetOnMessageCallback(server.onMessage)
	serverConn.SetOnErrorCallback(server.onError)
	serverConn.SetOnCloseCallback(server.onClose)
	serverConn.SetOnPacketRecvCallback(server.onPacket)

	return serverConn
}

func (server *ServerConnection) Start() {
	if server.GetOnConnectCallback() != nil {
		server.GetOnConnectCallback()(server)
	}

	// 每个连接启动3个线程，一个收，一个发，一个处理
	server.finish.Add(3)
	loopers := []func(Connection, *sync.WaitGroup){readLoop, writeLoop, handleServerLoop}
	for _, l := range loopers {
		looper := l // necessary
		go looper(server, server.finish)
	}
}

func (server *ServerConnection) Close() {
	server.once.Do(func() {
		if server.isClosed.CompareAndSet(false, true) {
			ok := server.GetOwner().connections.Remove(server.GetNetId())
			if !ok {
				glog.Errorf("ServerConnection conn %d %s remove failed, all size %d\n",
					server.GetNetId(), server.GetName(), server.GetOwner().connections.Size())
			}
			if server.GetOnCloseCallback() != nil {
				server.GetOnCloseCallback()(server)
			}

			close(server.GetCloseChannel())
			close(server.GetMessageSendChannel())
			close(server.GetHandlerReceiveChannel())
			close(server.GetPacketReceiveChannel())
			close(server.GetTimeOutChannel())

			// wait for all loops to finish
			server.finish.Wait()
			server.GetRawConn().Close()
			for _, id := range server.GetPendingTimers() {
				server.CancelTimer(id)
			}
			server.GetOwner().finish.Done()
		}
	})
}

func (server *ServerConnection) SetNetId(netid int64) {
	server.netid = netid
}

func (server *ServerConnection) GetNetId() int64 {
	return server.netid
}

func (server *ServerConnection) SetName(name string) {
	server.name = name
}

func (server *ServerConnection) GetName() string {
	return server.name
}

func (server *ServerConnection) Write(message Message) error {
	return asyncWrite(server, message)
}

func (server *ServerConnection) GetRawConn() net.Conn {
	return server.conn
}

func (server *ServerConnection) GetRemoteAddress() net.Addr {
	return server.conn.RemoteAddr()
}

func (server *ServerConnection) GetOwner() *TCPServer {
	return server.owner
}

func (server *ServerConnection) IsClosed() bool {
	return server.isClosed.Get()
}

func (server *ServerConnection) SetHeartBeat(beat int64) {
	server.heartBeat = beat
}

func (server *ServerConnection) GetHeartBeat() int64 {
	return server.heartBeat
}

func (server *ServerConnection) SetExtraData(extra interface{}) {
	server.extraData = extra
}

func (server *ServerConnection) GetExtraData() interface{} {
	return server.extraData
}

func (server *ServerConnection) SetMessageCodec(codec Codec) {
	server.messageCodec = codec
}

func (server *ServerConnection) GetMessageCodec() Codec {
	return server.messageCodec
}

func (server *ServerConnection) SetPendingTimers(pending []int64) {
	server.pendingTimers = pending
}

func (server *ServerConnection) GetPendingTimers() []int64 {
	return server.pendingTimers
}

func (server *ServerConnection) RunAt(timestamp time.Time, callback func(time.Time, interface{})) int64 {
	return runAt(server, timestamp, callback)
}

func (server *ServerConnection) RunAfter(duration time.Duration, callback func(time.Time, interface{})) int64 {
	return runAfter(server, duration, callback)
}

func (server *ServerConnection) RunEvery(interval time.Duration, callback func(time.Time, interface{})) int64 {
	return runEvery(server, interval, callback)
}

func (server *ServerConnection) GetTimingWheel() *TimingWheel {
	return server.GetOwner().GetTimingWheel()
}

func (server *ServerConnection) CancelTimer(timerId int64) {
	server.GetTimingWheel().CancelTimer(timerId)
}

func (server *ServerConnection) GetMessageSendChannel() chan []byte {
	return server.messageSendChan
}

func (server *ServerConnection) GetHandlerReceiveChannel() chan MessageHandler {
	return server.handlerRecvChan
}

func (server *ServerConnection) GetPacketReceiveChannel() chan *pool.Buffer {
	return server.packetRecvChan
}

func (server *ServerConnection) GetCloseChannel() chan struct{} {
	return server.closeConnChan
}

func (server *ServerConnection) GetTimeOutChannel() chan *OnTimeOut {
	return server.timeOutChan
}

func (server *ServerConnection) SetOnConnectCallback(callback func(Connection) bool) {
	server.onConnect = onConnectFunc(callback)
}

func (server *ServerConnection) GetOnConnectCallback() onConnectFunc {
	return server.onConnect
}

func (server *ServerConnection) SetOnMessageCallback(callback func(Message, Connection)) {
	server.onMessage = onMessageFunc(callback)
}

func (server *ServerConnection) GetOnMessageCallback() onMessageFunc {
	return server.onMessage
}

func (server *ServerConnection) SetOnErrorCallback(callback func()) {
	server.onError = onErrorFunc(callback)
}

func (server *ServerConnection) GetOnErrorCallback() onErrorFunc {
	return server.onError
}

func (server *ServerConnection) SetOnCloseCallback(callback func(Connection)) {
	server.onClose = onCloseFunc(callback)
}

func (server *ServerConnection) GetOnCloseCallback() onCloseFunc {
	return server.onClose
}

func (server *ServerConnection) SetOnPacketRecvCallback(callback onPacketRecvFunc) {
	server.onPacket = onPacketRecvFunc(callback)
}

func (server *ServerConnection) GetOnPacketRecvCallback() onPacketRecvFunc {
	return server.onPacket
}
