package net

import (
	"net"
	"os"
	"sync"
	"time"

	"github.com/golang/glog"

	. "smartgo/libs/utils"
)

type TCPServer struct {
	isRunning     *AtomicBoolean
	connections   *ConcurrentMap
	timingWheel   *TimingWheel
	workerPool    *WorkerPool
	finish        *sync.WaitGroup
	address       string
	closeServChan chan struct{}

	onConnect onConnectFunc
	onMessage onMessageFunc
	onClose   onCloseFunc
	onError   onErrorFunc
	onPacket  onPacketRecvFunc

	duration   time.Duration
	onSchedule onScheduleFunc
}

func NewTCPServer(addr string) Server {
	return &TCPServer{
		isRunning:     NewAtomicBoolean(true),
		connections:   NewConcurrentMap(),
		timingWheel:   NewTimingWheel(),
		workerPool:    NewWorkerPool(WORKERS),
		finish:        &sync.WaitGroup{},
		address:       addr,
		closeServChan: make(chan struct{}),
	}
}

func (server *TCPServer) IsRunning() bool {
	return server.isRunning.Get()
}

func (server *TCPServer) GetAllConnections() *ConcurrentMap {
	return server.connections
}

func (server *TCPServer) GetTimingWheel() *TimingWheel {
	return server.timingWheel
}

func (server *TCPServer) GetWorkerPool() *WorkerPool {
	return server.workerPool
}

func (server *TCPServer) GetServerAddress() string {
	return server.address
}

func (server *TCPServer) Start() {
	server.finish.Add(1)
	go server.timeOutLoop()

	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		glog.Fatalln(err)
	}
	defer listener.Close()

	for server.IsRunning() {
		conn, err := listener.Accept()
		if err != nil {
			glog.Errorln("Accept error - ", err)
			continue
		}

		conn = tlsWrapper(conn) // wrap as a tls connection if configured

		//config tcp conn
		tc, ok := conn.(*net.TCPConn)
		if ok {
			tc.SetKeepAlive(true)
			tc.SetKeepAlivePeriod(2 * time.Minute)
		}
		//conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		//conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

		/* Create a TCP connection upon accepting a new client, assign an net id
		   to it, then manage it in connections map, and start it */
		netid := netIdentifier.GetAndIncrement()
		tcpConn := NewServerConnection(netid, server, conn)
		tcpConn.SetName(tcpConn.GetRemoteAddress().String())
		if server.connections.Size() < MAX_CONNECTIONS {
			duration, onSchedule := server.GetOnScheduleCallback()
			if onSchedule != nil {
				tcpConn.RunEvery(duration, onSchedule)
			}
			server.connections.Put(netid, tcpConn)

			// put tcpConn.Start() run in another WaitGroup-synchronized go-routine
			server.finish.Add(1)
			go func() {
				tcpConn.Start()
			}()

			glog.Infof("Accepting client %s, net id %d, now total conn %d\n", tcpConn.GetName(), netid, server.connections.Size())
			//for v := range server.connections.IterValues() {
			//  glog.Infof("Client %s %t\n", v.(Connection).GetName(), v.(Connection).IsClosed())
			//}
		} else {
			glog.Warningf("MEET MAX CONNS LIMIT %d, refuse!\n", MAX_CONNECTIONS)
			tcpConn.Close()
		}
	}
}

// wait until all connections closed
func (server *TCPServer) Close() {
	if server.isRunning.CompareAndSet(true, false) {
		glog.Infoln("Server is to close --------- to close all client connections")
		for v := range server.GetAllConnections().IterValues() {
			c := v.(Connection)
			c.Close()
		}
		close(server.closeServChan)
		glog.Infoln("Server is to close --------- wait all sub routines to finish")
		server.finish.Wait()
		server.GetTimingWheel().Stop()
		glog.Infoln("Server is to close --------- final close, exit!")
		os.Exit(0)
	}
}

/* Retrieve the extra data(i.e. net id), and then redispatch
timeout callbacks to corresponding client connection, this
prevents one client from running callbacks of other clients */

func (server *TCPServer) timeOutLoop() {
	defer server.finish.Done()

	for {
		select {
		case <-server.closeServChan:
			return

		case timeout := <-server.GetTimingWheel().GetTimeOutChannel():
			netid := timeout.ExtraData.(int64)
			if conn, ok := server.connections.Get(netid); ok {
				tcpConn := conn.(Connection)
				if !tcpConn.IsClosed() {
					tcpConn.GetTimeOutChannel() <- timeout
				}
			} else {
				glog.Warningf("Invalid client %d\n", netid)
			}
		}
	}
}

func (server *TCPServer) GetTimeOutChannel() chan *OnTimeOut {
	return server.timingWheel.GetTimeOutChannel()
}

func (server *TCPServer) SetOnScheduleCallback(duration time.Duration, callback func(time.Time, interface{})) {
	server.duration = duration
	server.onSchedule = onScheduleFunc(callback)
}

func (server *TCPServer) GetOnScheduleCallback() (time.Duration, onScheduleFunc) {
	return server.duration, server.onSchedule
}

func (server *TCPServer) SetOnConnectCallback(callback func(Connection) bool) {
	server.onConnect = onConnectFunc(callback)
}

func (server *TCPServer) GetOnConnectCallback() onConnectFunc {
	return server.onConnect
}

func (server *TCPServer) SetOnMessageCallback(callback func(Message, Connection)) {
	server.onMessage = onMessageFunc(callback)
}

func (server *TCPServer) GetOnMessageCallback() onMessageFunc {
	return server.onMessage
}

func (server *TCPServer) SetOnCloseCallback(callback func(Connection)) {
	server.onClose = onCloseFunc(callback)
}

func (server *TCPServer) GetOnCloseCallback() onCloseFunc {
	return server.onClose
}

func (server *TCPServer) SetOnErrorCallback(callback func()) {
	server.onError = onErrorFunc(callback)
}

func (server *TCPServer) GetOnErrorCallback() onErrorFunc {
	return server.onError
}

func (server *TCPServer) SetOnPacketRecvCallback(callback onPacketRecvFunc) {
	server.onPacket = onPacketRecvFunc(callback)
}

func (server *TCPServer) GetOnPacketRecvCallback() onPacketRecvFunc {
	return server.onPacket
}
