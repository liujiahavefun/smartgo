package echo

import (
    "fmt"
    "testing"

    sessproto "smartgo/proto/sessevent"
    testproto "smartgo/proto/test"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
)

var signal *test.SignalTester

func server() {
    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start("127.0.0.1:7201")

    socket.RegisterMessage(server, "test.TestEchoACK", func(content interface{}, ses socket.Session) {
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

    socket.RegisterMessage(connector, "test.TestEchoACK", func(content interface{}, ses socket.Session) {
        msg := content.(*testproto.TestEchoACK)
        fmt.Println("client recv:", msg.String())
        signal.Done(1)
    })

    socket.RegisterMessage(connector, "sessevent.SessionConnected", func(content interface{}, ses socket.Session) {
        ses.Send(&testproto.TestEchoACK{
            Content: "hello",
        })
    })

    socket.RegisterMessage(connector, "sessevent.SessionConnectFailed", func(content interface{}, ses socket.Session) {
        msg := content.(*sessproto.SessionConnectFailed)
        fmt.Println(msg.Reason)
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "not recv data")
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

func TestEcho(t *testing.T) {
    signal = test.NewSignalTester(t)

    sessprotoRegisterMessage()
    testprotoRegisterMessage()
    socket.Init()

    server()
    client()
}