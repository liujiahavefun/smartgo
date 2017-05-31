package sendclose

import (
    "fmt"
    "testing"

    sessproto "smartgo/proto/sessevent"
    testproto "smartgo/proto/test"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
)

var signal *test.SignalTester

func runServer() {
    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(server, "sessevent.SessionConnected", func(content interface{}, ses socket.Session) {
        fmt.Println("recv SessionConnected, from peer: ", ses.FromPeer().Name())
        ses.FromPeer().SetName(ses.FromPeer().Name() + "_liujia")
    })

    socket.RegisterMessage(server, "test.TestEchoACK", func(content interface{}, ses socket.Session) {
        fmt.Println("recv TestEchoACK, from peer: ", ses.FromPeer().Name())

        msg := content.(*testproto.TestEchoACK)
        // 发包后关闭
        ses.Send(&testproto.TestEchoACK{
            Content: msg.Content,
        })

        if msg.Content != "noclose" {
            ses.Close()
        }
    })

    queue.StartLoop()
}

//客户端连接上后, 主动断开连接, 确保连接正常关闭
func testConnActiveClose() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(connector, "sessevent.SessionConnected", func(content interface{}, ses socket.Session) {
        signal.Done(1)
        //连接上发包,告诉服务器不要断开
        ses.Send(&testproto.TestEchoACK{
            Content: "noclose",
        })
    })

    socket.RegisterMessage(connector, "test.TestEchoACK", func(content interface{}, ses socket.Session) {
        msg := content.(*testproto.TestEchoACK)
        fmt.Println("client recv:", msg.String())
        signal.Done(2)

        // 客户端主动断开
        ses.Close()
    })

    socket.RegisterMessage(connector, "sessevent.SessionClosed", func(content interface{}, ses socket.Session) {
        msg := content.(*sessproto.SessionClosed)
        fmt.Println("close ok!", msg.Reason)

        // 正常断开
        signal.Done(3)
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "TestConnActiveClose not connected")
    signal.WaitAndExpect(2, "TestConnActiveClose not recv msg")
    signal.WaitAndExpect(3, "TestConnActiveClose not close")
}

// 接收封包后被断开
func testRecvDisconnected() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(connector, "sessevent.SessionConnected", func(content interface{}, ses socket.Session) {
        // 连接上发包
        ses.Send(&testproto.TestEchoACK{
            Content: "data",
        })

        signal.Done(1)
    })

    socket.RegisterMessage(connector, "test.TestEchoACK", func(content interface{}, ses socket.Session) {
        msg := content.(*testproto.TestEchoACK)
        fmt.Println("client recv:", msg.String())
        signal.Done(2)
    })

    socket.RegisterMessage(connector, "sessevent.SessionClosed", func(content interface{}, ses socket.Session) {
        // 断开
        signal.Done(3)
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "TestRecvDisconnected not connected")
    signal.WaitAndExpect(2, "TestRecvDisconnected not recv msg")
    signal.WaitAndExpect(3, "TestRecvDisconnected not closed")
}

func sessprotoRegisterMessage() {
    fmt.Println("sessevent RegisterMessage")
    // session.proto
    socket.RegisterMessageMeta("sessevent.SessionAccepted", (*sessproto.SessionAccepted)(nil), 2136350511)
    socket.RegisterMessageMeta("sessevent.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1213847952)
    socket.RegisterMessageMeta("sessevent.SessionConnected", (*sessproto.SessionConnected)(nil), 4228538224)
    socket.RegisterMessageMeta("sessevent.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1278926828)
    socket.RegisterMessageMeta("sessevent.SessionClosed", (*sessproto.SessionClosed)(nil), 2830250790)
    socket.RegisterMessageMeta("sessevent.SessionError", (*sessproto.SessionError)(nil), 3227768243)

}

func testprotoRegisterMessage() {
    fmt.Println("test RegisterMessage")
    // test.proto
    socket.RegisterMessageMeta("test.TestEchoACK", (*testproto.TestEchoACK)(nil), 509149489)
}

func TestClose(t *testing.T) {
    signal = test.NewSignalTester(t)

    sessprotoRegisterMessage()
    testprotoRegisterMessage()
    socket.Init()

    runServer()

    testConnActiveClose()
    testRecvDisconnected()
}