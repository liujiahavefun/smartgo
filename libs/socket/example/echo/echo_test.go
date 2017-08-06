package echo

import (
    "fmt"
    "testing"

    sessproto "smartgo/proto/session_event"
    testproto "smartgo/proto/test_event"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
)

var signal *test.SignalTester

func server() {
    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(server, "test_event.TestEchoACK", func(content interface{}, ses socket.Session) {
        msg := content.(*testproto.TestEchoACK)
        fmt.Println("server recv:", msg.String())
        ses.Send(&testproto.TestEchoACK{
            Content: msg.String(),
        })
    })

    queue.StartLoop()
}

func client() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(connector, "test_event.TestEchoACK", func(content interface{}, ses socket.Session) {
        msg := content.(*testproto.TestEchoACK)
        fmt.Println("client recv:", msg.String())
        signal.Done(1)
    })

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, ses socket.Session) {
        ses.Send(&testproto.TestEchoACK{
            Content: "hello",
        })
    })

    socket.RegisterMessage(connector, "session_event.SessionConnectFailed", func(content interface{}, ses socket.Session) {
        msg := content.(*sessproto.SessionConnectFailed)
        fmt.Println(msg.Reason)
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "not recv data")
}

func sessprotoRegisterMessage() {
    fmt.Println("sessevent RegisterMessage")
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
    socket.RegisterMessageMeta("test_event.TestEchoACK", (*testproto.TestEchoACK)(nil), 509149489)
}

func TestEcho(t *testing.T) {
    signal = test.NewSignalTester(t)

    sessprotoRegisterMessage()
    testprotoRegisterMessage()
    socket.Init()

    server()
    client()
}