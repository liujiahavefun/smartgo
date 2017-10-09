package sendclose

import (
    "fmt"
    "testing"

    sessproto "smartgo/proto/session_event"
    testproto "smartgo/proto/test_event"
    "smartgo/libs/utils"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
    "strings"
    "time"
)

var signal *test.SignalTester

func runServer() {
    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(server, "session_event.SessionConnected", func(content interface{}, ses socket.Session) {
        fmt.Println("Server: recv SessionConnected from peer", ses.FromPeer().Name(), utils.GoID())
        ses.FromPeer().SetName(ses.FromPeer().Name() + "_liujia")
    })

    socket.RegisterMessage(server, "test_event.TestEchoACK", func(content interface{}, ses socket.Session) {
        fmt.Println("Server: recv TestEchoACK from peer", ses.FromPeer().Name(), utils.GoID())

        msg := content.(*testproto.TestEchoACK)
        // 发包后关闭
        ses.Send(&testproto.TestEchoACK{
            Content: msg.Content + "_liujia",
        })

        if msg.Content != "noclose" {
            fmt.Println("Server: close session to peer ", ses.FromPeer().Name(), utils.GoID())
            ses.Close()
            signal.Done(2)
        }
    })

    queue.StartLoop()
}

//客户端连接上后, 主动断开连接, 确保连接正常关闭
func testConnActiveClose() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, ses socket.Session) {
        signal.Done(1)
        //连接上发包,告诉服务器不要断开
        ses.Send(&testproto.TestEchoACK{
            Content: "noclose",
        })
    })

    socket.RegisterMessage(connector, "test_event.TestEchoACK", func(content interface{}, ses socket.Session) {
        msg := content.(*testproto.TestEchoACK)
        fmt.Println("Client: recv ", msg.String())

        if strings.Contains(msg.Content, "liujia") == true {
            queue.PostDelayed(connector, 5*time.Second, func() {
                fmt.Println("Client: postDelayed ", utils.GoID())
                //告诉服务器断开
                ses.Send(&testproto.TestEchoACK{
                    Content: "close",
                })
            })
            return
        }

        //fmt.Println("Server: close session to peer ", ses.FromPeer().Name())
        //signal.Done(2)
        //客户端主动断开
        //ses.Close()
    })

    socket.RegisterMessage(connector, "session_event.SessionClosed", func(content interface{}, ses socket.Session) {
        msg := content.(*sessproto.SessionClosed)
        fmt.Println("Client: recv SessionClosed", msg.Reason, utils.GoID())

        // 正常断开
        signal.Done(3)
    })

    socket.NewTimer(queue, 1*time.Second, func(t *socket.Timer) {
        fmt.Println("Client: timer", utils.GoID())
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

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, ses socket.Session) {
        // 连接上发包
        ses.Send(&testproto.TestEchoACK{
            Content: "data",
        })

        signal.Done(1)
    })

    socket.RegisterMessage(connector, "test_event.TestEchoACK", func(content interface{}, ses socket.Session) {
        msg := content.(*testproto.TestEchoACK)
        fmt.Println("client recv:", msg.String())
        signal.Done(2)
    })

    socket.RegisterMessage(connector, "session_event.SessionClosed", func(content interface{}, ses socket.Session) {
        // 断开
        signal.Done(3)
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "TestRecvDisconnected not connected")
    signal.WaitAndExpect(2, "TestRecvDisconnected not recv msg")
    signal.WaitAndExpect(3, "TestRecvDisconnected not closed")
}

func sessprotoRegisterMessage() {
    fmt.Println("session_event RegisterMessage")
    // session.proto
    socket.RegisterMessageMeta("session_event.SessionAccepted", (*sessproto.SessionAccepted)(nil), 348117910)
    socket.RegisterMessageMeta("session_event.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1978788392)
    socket.RegisterMessageMeta("session_event.SessionConnected", (*sessproto.SessionConnected)(nil), 3543838007)
    socket.RegisterMessageMeta("session_event.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1720533237)
    socket.RegisterMessageMeta("session_event.SessionClosed", (*sessproto.SessionClosed)(nil), 90181607)
    socket.RegisterMessageMeta("session_event.SessionError", (*sessproto.SessionError)(nil), 1937281175)
}

func testprotoRegisterMessage() {
    fmt.Println("test_event RegisterMessage")
    // test.proto
    socket.RegisterMessageMeta("test_event.TestEchoACK", (*testproto.TestEchoACK)(nil), 524159734)
}

func TestClose(t *testing.T) {
    signal = test.NewSignalTesterTimeout(t, 10)

    sessprotoRegisterMessage()
    testprotoRegisterMessage()
    socket.Init()

    runServer()

    testConnActiveClose()
    //testRecvDisconnected()
}