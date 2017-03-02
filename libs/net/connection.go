package net

import (
	"net"
	"sync"
	"time"

	"github.com/golang/glog"

	"smartgo/libs/pool"
	. "smartgo/libs/utils"
)

/*
* server interface for TCP/TLS-TCP server
 */
type Connection interface {
	Start()
	Close()
	IsClosed() bool
	Write(message Message) error

	SetNetId(netid int64)
	GetNetId() int64

	SetName(name string)
	GetName() string

	GetRawConn() net.Conn
	GetRemoteAddress() net.Addr

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

	SetOnErrorCallback(callback func())
	GetOnErrorCallback() onErrorFunc

	SetOnCloseCallback(callback func(Connection))
	GetOnCloseCallback() onCloseFunc

	SetOnMessageCallback(callback func(Message, Connection))
	GetOnMessageCallback() onMessageFunc

	SetOnPacketRecvCallback(callback onPacketRecvFunc)
	GetOnPacketRecvCallback() onPacketRecvFunc

	RunAt(t time.Time, cb func(time.Time, interface{})) int64
	RunAfter(d time.Duration, cb func(time.Time, interface{})) int64
	RunEvery(i time.Duration, cb func(time.Time, interface{})) int64
	GetTimingWheel() *TimingWheel

	CancelTimer(timerId int64)

	GetMessageSendChannel() chan []byte
	GetHandlerReceiveChannel() chan MessageHandler
	GetPacketReceiveChannel() chan *pool.Buffer
	GetCloseChannel() chan struct{}
	GetTimeOutChannel() chan *OnTimeOut
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

/*
* runAt/runAfter/runEvery，都是用来注册一个连接上的定时任务的
*/
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