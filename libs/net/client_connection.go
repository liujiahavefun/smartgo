package net

import (
	"net"
	"sync"
	"time"

	"github.com/golang/glog"

	"smartgo/libs/pool"
	. "smartgo/libs/utils"
	"crypto/tls"
)

// implements Client Side Connection
type ClientConnection struct {
	netid         int64
	name          string
	address       string
	conn          net.Conn
	finish        *sync.WaitGroup
	isClosed      *AtomicBoolean
	reconnectable bool

	once          *sync.Once
	pendingTimers []int64
	timingWheel   *TimingWheel
	messageCodec  Codec

	heartBeat     int64
	extraData     interface{}

	messageSendChan chan []byte
	handlerRecvChan chan MessageHandler
	packetRecvChan  chan *pool.Buffer
	closeConnChan   chan struct{}

	onConnect onConnectFunc
	onMessage onMessageFunc
	onClose   onCloseFunc
	onError   onErrorFunc
	onPacket  onPacketRecvFunc
}

func NewTLSClientConnection(netid int64, address string, reconnectable bool, config *tls.Config, onPacket onPacketRecvFunc) Connection {
	c, err := tls.Dial("tcp", address, config)
	if err != nil {
		glog.Fatalln("NewTLSClientConnection", err)
	}
	return &ClientConnection{
		netid:           netid,
		name:            c.RemoteAddr().String(),
		address:         address,
		heartBeat:       time.Now().UnixNano(),
		isClosed:        NewAtomicBoolean(false),
		once:            &sync.Once{},
		pendingTimers:   []int64{},
		timingWheel:     NewTimingWheel(),
		conn:            c,
		messageCodec:    NewFixedLengthHeaderCodec(HEADER_SIZE, MIN_PACKET_SIZE, MAX_PACKET_SIZE, GLOBAL_BINARY_BYTE_ORDER),
		finish:          &sync.WaitGroup{},
		reconnectable:   reconnectable,
		messageSendChan: make(chan []byte, 1024),
		handlerRecvChan: make(chan MessageHandler, 1024),
		packetRecvChan:  make(chan *pool.Buffer, 1024),
		closeConnChan:   make(chan struct{}),
		onPacket:        onPacket,
	}
}

func NewClientConnection(netid int64, address string, reconnectable bool, onPacket onPacketRecvFunc) (Connection, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &ClientConnection{
		netid:           netid,
		name:            c.RemoteAddr().String(),
		address:         address,
		heartBeat:       time.Now().UnixNano(),
		isClosed:        NewAtomicBoolean(false),
		once:            &sync.Once{},
		pendingTimers:   []int64{},
		timingWheel:     NewTimingWheel(),
		conn:            c,
		messageCodec:    NewFixedLengthHeaderCodec(HEADER_SIZE, MIN_PACKET_SIZE, MAX_PACKET_SIZE, GLOBAL_BINARY_BYTE_ORDER),
		finish:          &sync.WaitGroup{},
		reconnectable:   reconnectable,
		messageSendChan: make(chan []byte, 1024),
		handlerRecvChan: make(chan MessageHandler, 1024),
		packetRecvChan:  make(chan *pool.Buffer, 1024),
		closeConnChan:   make(chan struct{}),
		onPacket:        onPacket,
	}, nil
}

func (client *ClientConnection) Start() {
	if client.GetOnConnectCallback() != nil {
		client.GetOnConnectCallback()(client)
	}

	client.finish.Add(3)
	loopers := []func(Connection, *sync.WaitGroup){readLoop, writeLoop, handleClientLoop}
	for _, l := range loopers {
		looper := l // necessary
		go looper(client, client.finish)
	}
}

func (client *ClientConnection) Close() {
	done := false
	client.once.Do(func() {
		if client.isClosed.CompareAndSet(false, true) {
			if client.GetOnCloseCallback() != nil {
				client.GetOnCloseCallback()(client)
			}

			close(client.GetCloseChannel())
			close(client.GetMessageSendChannel())
			close(client.GetHandlerReceiveChannel())
			close(client.GetPacketReceiveChannel())
			close(client.GetTimeOutChannel())
			client.GetTimingWheel().Stop()

			// close tcp conn
			client.GetRawConn().Close()

			// wait for all loops to finish
			client.finish.Wait()
			done = true
		}
	})

	if done && client.reconnectable {
		client.reconnect()
	}
}

func (client *ClientConnection) reconnect() {
	c, err := net.Dial("tcp", client.address)
	if err != nil {
		glog.Fatalln("ClientConnection", err)
	}
	client.name = c.RemoteAddr().String()
	client.heartBeat = time.Now().UnixNano()
	client.extraData = nil
	client.once = &sync.Once{}
	client.pendingTimers = []int64{}
	client.timingWheel = NewTimingWheel()
	client.conn = c
	client.messageSendChan = make(chan []byte, 1024)
	client.handlerRecvChan = make(chan MessageHandler, 1024)
	client.closeConnChan = make(chan struct{})
	client.Start()
	client.isClosed.CompareAndSet(true, false)
}

func (client *ClientConnection) SetNetId(netid int64) {
	client.netid = netid
}

func (client *ClientConnection) GetNetId() int64 {
	return client.netid
}

func (client *ClientConnection) SetName(name string) {
	client.name = name
}

func (client *ClientConnection) GetName() string {
	return client.name
}

func (client *ClientConnection) GetRawConn() net.Conn {
	return client.conn
}

func (client *ClientConnection) GetCloseChannel() chan struct{} {
	return client.closeConnChan
}

func (client *ClientConnection) GetRemoteAddress() net.Addr {
	return client.conn.RemoteAddr()
}

func (client *ClientConnection) IsClosed() bool {
	return client.isClosed.Get()
}

func (client *ClientConnection) Write(message Message) error {
	return asyncWrite(client, message)
}

func (client *ClientConnection) SetHeartBeat(beat int64) {
	client.heartBeat = beat
}

func (client *ClientConnection) GetHeartBeat() int64 {
	return client.heartBeat
}

func (client *ClientConnection) SetExtraData(extra interface{}) {
	client.extraData = extra
}

func (client *ClientConnection) GetExtraData() interface{} {
	return client.extraData
}

func (client *ClientConnection) SetMessageCodec(codec Codec) {
	client.messageCodec = codec
}

func (client *ClientConnection) GetMessageCodec() Codec {
	return client.messageCodec
}

func (client *ClientConnection) SetOnConnectCallback(callback func(Connection) bool) {
	client.onConnect = onConnectFunc(callback)
}

func (client *ClientConnection) GetOnConnectCallback() onConnectFunc {
	return client.onConnect
}

func (client *ClientConnection) SetOnMessageCallback(callback func(Message, Connection)) {
	client.onMessage = onMessageFunc(callback)
}

func (client *ClientConnection) GetOnMessageCallback() onMessageFunc {
	return client.onMessage
}

func (client *ClientConnection) SetOnErrorCallback(callback func()) {
	client.onError = onErrorFunc(callback)
}

func (client *ClientConnection) GetOnErrorCallback() onErrorFunc {
	return client.onError
}

func (client *ClientConnection) SetOnCloseCallback(callback func(Connection)) {
	client.onClose = onCloseFunc(callback)
}

func (client *ClientConnection) GetOnCloseCallback() onCloseFunc {
	return client.onClose
}

func (server *ClientConnection) SetOnPacketRecvCallback(callback onPacketRecvFunc) {
	server.onPacket = onPacketRecvFunc(callback)
}

func (server *ClientConnection) GetOnPacketRecvCallback() onPacketRecvFunc {
	return server.onPacket
}

func (client *ClientConnection) RunAt(timestamp time.Time, callback func(time.Time, interface{})) int64 {
	return runAt(client, timestamp, callback)
}

func (client *ClientConnection) RunAfter(duration time.Duration, callback func(time.Time, interface{})) int64 {
	return runAfter(client, duration, callback)
}

func (client *ClientConnection) RunEvery(interval time.Duration, callback func(time.Time, interface{})) int64 {
	return runEvery(client, interval, callback)
}

func (client *ClientConnection) GetTimingWheel() *TimingWheel {
	return client.timingWheel
}

func (client *ClientConnection) SetPendingTimers(pending []int64) {
	client.pendingTimers = pending
}

func (client *ClientConnection) GetPendingTimers() []int64 {
	return client.pendingTimers
}

func (client *ClientConnection) CancelTimer(timerId int64) {
	client.GetTimingWheel().CancelTimer(timerId)
}

func (client *ClientConnection) GetMessageSendChannel() chan []byte {
	return client.messageSendChan
}

func (client *ClientConnection) GetHandlerReceiveChannel() chan MessageHandler {
	return client.handlerRecvChan
}

func (client *ClientConnection) GetPacketReceiveChannel() chan *pool.Buffer {
	return client.packetRecvChan
}

func (client *ClientConnection) GetTimeOutChannel() chan *OnTimeOut {
	return client.timingWheel.GetTimeOutChannel()
}
