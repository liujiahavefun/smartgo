package net

import (
	"crypto/tls"
	"github.com/golang/glog"
	"net"
	"smartgo/libs/pool"
	. "smartgo/libs/utils"
	"sync"
	"time"
)

type Connection interface {
	SetNetId(netid int64)
	GetNetId() int64

	SetName(name string)
	GetName() string

	SetHeartBeat(beat int64)
	GetHeartBeat() int64

	SetExtraData(extra interface{})
	GetExtraData() interface{}

	SetMessageCodec(codec Codec)
	GetMessageCodec() Codec

	SetPendingTimers(pending []int64)
	GetPendingTimers() []int64

	SetOnConnectCallback(callback func(Connection) bool)
	GetOnConnectCallback() onConnectFunc
	SetOnMessageCallback(callback func(Message, Connection))
	GetOnMessageCallback() onMessageFunc
	SetOnErrorCallback(callback func())
	GetOnErrorCallback() onErrorFunc
	SetOnCloseCallback(callback func(Connection))
	GetOnCloseCallback() onCloseFunc

	SetOnPacketRecvCallback(callback onPacketRecvFunc)
	GetOnPacketRecvCallback() onPacketRecvFunc

	Start()
	Close()
	IsClosed() bool
	Write(message Message) error

	RunAt(t time.Time, cb func(time.Time, interface{})) int64
	RunAfter(d time.Duration, cb func(time.Time, interface{})) int64
	RunEvery(i time.Duration, cb func(time.Time, interface{})) int64
	GetTimingWheel() *TimingWheel

	CancelTimer(timerId int64)

	GetRawConn() net.Conn
	GetMessageSendChannel() chan []byte
	GetHandlerReceiveChannel() chan MessageHandler
	GetPacketReceiveChannel() chan *pool.Buffer
	GetCloseChannel() chan struct{}
	GetTimeOutChannel() chan *OnTimeOut

	GetRemoteAddress() net.Addr
}

// type Reconnectable interface{
//   Reconnect()
// }

type ServerSide interface {
	GetOwner() *TCPServer
}

type ServerSideConnection interface {
	Connection
	ServerSide
}

type ClientSideConnection interface {
	Connection
	// Reconnectable
}

// implements ServerSideConnection
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

func (server *ServerConnection) GetOwner() *TCPServer {
	return server.owner
}

func (server *ServerConnection) Start() {
	if server.GetOnConnectCallback() != nil {
		server.GetOnConnectCallback()(server)
	}

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

func (server *ServerConnection) IsClosed() bool {
	return server.isClosed.Get()
}

func (server *ServerConnection) SetPendingTimers(pending []int64) {
	server.pendingTimers = pending
}

func (server *ServerConnection) GetPendingTimers() []int64 {
	return server.pendingTimers
}

func (server *ServerConnection) Write(message Message) error {
	return asyncWrite(server, message)
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

func (server *ServerConnection) GetRawConn() net.Conn {
	return server.conn
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

func (server *ServerConnection) GetRemoteAddress() net.Addr {
	return server.conn.RemoteAddr()
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

// implements Connection
type ClientConnection struct {
	netid         int64
	name          string
	address       string
	heartBeat     int64
	extraData     interface{}
	isClosed      *AtomicBoolean
	once          *sync.Once
	pendingTimers []int64
	timingWheel   *TimingWheel
	conn          net.Conn
	messageCodec  Codec
	finish        *sync.WaitGroup
	reconnectable bool

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

func (client *ClientConnection) IsClosed() bool {
	return client.isClosed.Get()
}

func (client *ClientConnection) Write(message Message) error {
	return asyncWrite(client, message)
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

func (client *ClientConnection) GetRawConn() net.Conn {
	return client.conn
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

func (client *ClientConnection) GetCloseChannel() chan struct{} {
	return client.closeConnChan
}

func (client *ClientConnection) GetTimeOutChannel() chan *OnTimeOut {
	return client.timingWheel.GetTimeOutChannel()
}

func (client *ClientConnection) GetRemoteAddress() net.Addr {
	return client.conn.RemoteAddr()
}

func runAt(conn Connection, timestamp time.Time, callback func(time.Time, interface{})) int64 {
	timeout := NewOnTimeOut(conn.GetNetId(), callback)
	var id int64 = -1
	if conn.GetTimingWheel() != nil {
		id = conn.GetTimingWheel().AddTimer(timestamp, 0, timeout)
		if id >= 0 {
			pending := conn.GetPendingTimers()
			pending = append(pending, id)
			conn.SetPendingTimers(pending)
		}
	}
	return id
}

func runAfter(conn Connection, duration time.Duration, callback func(time.Time, interface{})) int64 {
	delay := time.Now().Add(duration)
	var id int64 = -1
	if conn.GetTimingWheel() != nil {
		id = conn.RunAt(delay, callback)
	}
	return id
}

func runEvery(conn Connection, interval time.Duration, callback func(time.Time, interface{})) int64 {
	delay := time.Now().Add(interval)
	timeout := NewOnTimeOut(conn.GetNetId(), callback)
	var id int64 = -1
	if conn.GetTimingWheel() != nil {
		id = conn.GetTimingWheel().AddTimer(delay, interval, timeout)
		if id >= 0 {
			pending := conn.GetPendingTimers()
			pending = append(pending, id)
			conn.SetPendingTimers(pending)
		}
	}
	return id
}

func asyncWrite(conn Connection, message Message) error {
	if conn.IsClosed() {
		return ErrorConnClosed
	}

	packet, err := conn.GetMessageCodec().Encode(message)
	if err != nil {
		glog.Errorln("asyncWrite", err)
		return err
	}

	select {
	case conn.GetMessageSendChannel() <- packet:
		return nil
	default:
		return ErrorWouldBlock
	}
}

//check if the connection is closed in NON-BLOCK way
func isClosed(conn Connection) bool {
	if conn.IsClosed() {
		return true
	}

	select {
	case <-conn.GetCloseChannel():
		return true
	default:
		return false
	}
}

/* readLoop() blocking read from connection, deserialize bytes into message,
then find corresponding handler, put it into channel */
func readLoop(conn Connection, finish *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			glog.Errorf("readLoop panics: %v\n", p)
		}
		finish.Done()
		conn.Close()
		glog.Errorf("readLoop exit")
	}()

	for {
		/*
			select {
			case <-conn.GetCloseChannel():
				return

			default:
			}
		*/
		if isClosed(conn) {
			glog.Errorf("readLoop: recevie close signal")
			return
		}

		//TODO: 某些错误可以忽略，比如包太大之类的
		buf, err := conn.GetMessageCodec().Decode(conn)
		if err != nil {
			/*
				if _, ok := err.(*net.OpError); ok {
					continue
				}
			*/

			if err == ErrorkrPacketTooLarge || err == ErrorkrPacketInvalidBody {
				conn.SetHeartBeat(time.Now().UnixNano())
				continue
			}

			// 读写超时算错误么？
			netError, ok := err.(net.Error)
			if ok && netError.Timeout() {
				continue
			}

			buf.Free()

			glog.Errorf("readLoop Error decoding message - %s\n", err)
			if _, ok := err.(ErrorUndefined); ok {
				conn.SetHeartBeat(time.Now().UnixNano())
				continue
			}
			return
		}

		// update heart beat timestamp
		conn.SetHeartBeat(time.Now().UnixNano())

		if !conn.IsClosed() {
			glog.Errorf("delivery msg to packet recv channel\n")
			conn.GetPacketReceiveChannel() <- buf
		}

		/*
			handlerFactory := HandlerMap.Get(msg.MessageNumber())
			if handlerFactory == nil {
				if conn.GetOnMessageCallback() != nil {
					glog.Infof("readLoop Message %d call onMessage()\n", msg.MessageNumber())
					conn.GetOnMessageCallback()(msg, conn)
				} else {
					glog.Infof("readLoop No handler or onMessage() found for message %d", msg.MessageNumber())
				}
				continue
			}

			// send handler to handleLoop
			handler := handlerFactory(conn.GetNetId(), msg)
			if !conn.IsClosed() {
				conn.GetHandlerReceiveChannel() <- handler
			}
		*/
	}
}

/* writeLoop() receive message from channel, serialize it into bytes,
then blocking write into connection */
func writeLoop(conn Connection, finish *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			glog.Errorf("writeLoop panics: %v", p)
		}
		for packet := range conn.GetMessageSendChannel() {
			if packet != nil {
				if _, err := conn.GetRawConn().Write(packet); err != nil {
					glog.Errorf("writeLoop Error writing data - %s\n", err)
				}
			}
		}
		finish.Done()
		conn.Close()
		glog.Errorf("writeLoop exit")
	}()

	for {
		select {
		case <-conn.GetCloseChannel():
			glog.Errorf("writeLoop: recevie close signal")
			return

		case packet := <-conn.GetMessageSendChannel():
			if packet != nil {
				if _, err := conn.GetRawConn().Write(packet); err != nil {
					glog.Errorf("writeLoop Error writing data - %s\n", err)
					return
				}
			}
		}
	}
}

// handleServerLoop() - put handler or timeout callback into worker go-routines
func handleServerLoop(conn Connection, finish *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			glog.Errorf("handleServerLoop panics: %v", p)
		}
		finish.Done()
		conn.Close()
		glog.Errorf("handleServerLoop exit")
	}()

	for {
		select {
		case <-conn.GetCloseChannel():
			glog.Errorf("handleServerLoop: recevie client close signal")
			return

		/*
			case handler := <-conn.GetHandlerReceiveChannel():
				if !IsNil(handler) {
					serverConn, ok := conn.(*ServerConnection)
					if ok {
						serverConn.GetOwner().workerPool.Put(conn.GetNetId(), func() {
							handler.Process(conn)
						})
					}
				}
		*/
		case buf := <-conn.GetPacketReceiveChannel():
			if buf != nil {
				serverConn, ok := conn.(*ServerConnection)
				if ok {
					onPacket := serverConn.GetOnPacketRecvCallback()
					if onPacket != nil {
						serverConn.GetOwner().workerPool.Put(conn.GetNetId(), func() {
							defer func() {
								if p := recover(); p != nil {
									glog.Errorf("HandlePacket panics: %v", p)
									conn.Close()
								}
							}()

							defer buf.Free()

							handler, ok := onPacket(conn, buf)
							if handler != nil && ok {
								handler()
							}
						})
					}
				}
			}

		case timeout := <-conn.GetTimeOutChannel():
			if timeout != nil {
				extraData := timeout.ExtraData.(int64)
				if extraData != conn.GetNetId() {
					glog.Warningf("handleServerLoop time out of %d running on client %d", extraData, conn.GetNetId())
				}
				serverConn, ok := conn.(*ServerConnection)
				if ok {
					serverConn.GetOwner().workerPool.Put(conn.GetNetId(), func() {
						timeout.Callback(time.Now(), conn)
					})
				} else {
					glog.Errorf("handleServerLoop conn %s is not of type *ServerConnection\n", conn.GetName())
				}
			}
		}
	}
}

// handleClientLoop() - run handler or timeout callback in handleLoop() go-routine
func handleClientLoop(conn Connection, finish *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			glog.Errorf("handleClientLoop panics: %v", p)
		}
		finish.Done()
		conn.Close()
		glog.Errorf("handleClientLoop exit")
	}()

	for {
		select {
		case <-conn.GetCloseChannel():
			glog.Errorf("handleClientLoop: recevie close signal")
			return

		/*
			case handler := <-conn.GetHandlerReceiveChannel():
				if !IsNil(handler) {
					handler.Process(conn)
				}
		*/
		case buf := <-conn.GetPacketReceiveChannel():
			if buf != nil {
				clientConn, ok := conn.(*ClientConnection)
				if ok {
					onPacket := clientConn.GetOnPacketRecvCallback()
					if onPacket != nil {
						func() {
							defer func() {
								if p := recover(); p != nil {
									glog.Errorf("HandlePacket panics: %v", p)
									conn.Close()
								}
							}()
							defer buf.Free()

							handler, ok := onPacket(conn, buf)
							if handler != nil && ok {
								handler()
							}
						}()
					}
				}
			}

		case timeout := <-conn.GetTimeOutChannel():
			if timeout != nil {
				extraData := timeout.ExtraData.(int64)
				if extraData != conn.GetNetId() {
					glog.Warningf("handleClientLoop time out of %d running on client %d", extraData, conn.GetNetId())
				}
				timeout.Callback(time.Now(), conn)
			}
		}
	}
}
