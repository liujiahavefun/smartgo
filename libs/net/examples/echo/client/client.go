package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"smartgo/libs/net"
	"smartgo/libs/net/examples/echo"
	"smartgo/libs/pool"
	"time"
)

func init() {
	globalPool = pool.NewBufferPool(128, 128)
	defaultHandler = func() {}
}

var (
	globalPool     *pool.BufferPool
	defaultHandler func()
)

func HandlePacket(conn net.Connection, packet *pool.Buffer) (net.HandlerProc, bool) {
	glog.Errorf("client handle packet: %v\n", packet.Data)
	packet.ResetSeeker()
	major := packet.ReadUint16BE()
	minor := packet.ReadUint16BE()

	if major == 1 && minor == 1 {
		s := packet.ReadString()
		glog.Errorf("Receive Server Message: ", s)
	}

	return nil, false
}

func main() {
	flag.Parse()

	tcpConnection := net.NewClientConnection(0, "123.56.88.196:9100", false, HandlePacket)

	tcpConnection.SetOnConnectCallback(func(client net.Connection) bool {
		glog.Errorf("On connect")
		return true
	})

	tcpConnection.SetOnErrorCallback(func() {
		glog.Errorf("On error")
	})

	tcpConnection.SetOnCloseCallback(func(client net.Connection) {
		glog.Errorf("On close")
	})

	tcpConnection.SetOnMessageCallback(func(msg net.Message, c net.Connection) {
		echoMessage, ok := msg.(echo.EchoMessage)
		if ok {
			fmt.Printf("received echo msg: %v\n", echoMessage.EchoString)
		} else {
			fmt.Printf("recevied msg\n")
		}
	})

	echoMessage := echo.EchoMessage{
		EchoString: "hello, world",
	}

	tcpConnection.RunAt(time.Now().Add(time.Second*20), func(now time.Time, data interface{}) {
		client := data.(net.Connection)
		glog.Errorf("Closing after 20 seconds")
		client.Close()
	})

	tcpConnection.Start()

	for i := 0; i < 2; i++ {
		err := tcpConnection.Write(echoMessage)
		if err != nil {
			glog.Errorln("Write Error: ", err)
		}
		time.Sleep(time.Second)
	}

	time.Sleep(1 * time.Second)
	tcpConnection.Close()
}
