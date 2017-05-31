package main

import (
    "fmt"

    sessproto "smartgo/proto/sessevent"
    testproto "smartgo/proto/test"
    "smartgo/libs/socket"
)

func init() {
    sessprotoRegisterMessage()
    testprotoRegisterMessage()
}

func startServer(address string) {
    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start(address)

    socket.RegisterMessage(server, "sessevent.SessionAccepted", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionAccepted)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "sessevent.SessionAcceptFailed", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionAcceptFailed)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "sessevent.SessionConnected", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionConnected)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "sessevent.SessionConnectFailed", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionConnectFailed)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "sessevent.SessionError", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionError)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "sessevent.SessionClosed", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionClosed)
        fmt.Println(msg)
    })

    queue.StartLoop()
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