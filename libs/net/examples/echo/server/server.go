package main

import (
	"flag"
	"github.com/golang/glog"
	"runtime"
	"smartgo/libs/net"
	"smartgo/libs/net/examples/echo"
	"smartgo/libs/pool"
)

func init() {
	defaultHandler = func() {}
}

var (
	defaultHandler func()
)

type EchoServer struct {
	net.Server
}

func NewEchoServer(addr string) *EchoServer {
	return &EchoServer{
		net.NewTCPServer(addr),
	}
}

func HandlePacket(conn net.Connection, packet *pool.Buffer) (net.HandlerProc, bool) {
	glog.Errorf("server handle packet: %v\n", packet.Data)
	packet.ResetSeeker()
	major := packet.ReadUint16BE()
	minor := packet.ReadUint16BE()

	if major == 1 && minor == 1 {
		s := packet.ReadString()
		glog.Errorf("read string: %v\n", s)
		em := echo.EchoMessage{
			EchoString: ("朕知道了: " + s),
		}

		return func() {
			conn.Write(em)
		}, true
	}

	return nil, false
}

func main() {
	flag.Set("v", "5")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	echoServer := NewEchoServer(":9100")
	defer echoServer.Close()

	echoServer.SetOnConnectCallback(func(client net.Connection) bool {
		glog.Infoln("On connect")
		return true
	})

	echoServer.SetOnErrorCallback(func() {
		glog.Infoln("On error")
	})

	echoServer.SetOnCloseCallback(func(client net.Connection) {
		glog.Infoln("Closing client")
	})

	//echoServer.SetOnMessageCallback(func(msg net.Message, client net.Connection) {
	//	glog.Infoln("Receving message")
	//})

	echoServer.SetOnPacketRecvCallback(HandlePacket)

	echoServer.Start()
}
