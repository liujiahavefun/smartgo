package main

import (
    "fmt"

    sessproto "smartgo/proto/session_event"
    loginproto "smartgo/proto/login_event"
    "smartgo/libs/socket"
)

func init() {
    sessionprotoRegisterMessage()
    loginprotoRegisterMessage()
}

func startServer(address string) {
    queue := socket.NewEventQueue()
    server := socket.NewTcpServer(queue).Start(address)

    socket.RegisterMessage(server, "session_event.SessionAccepted", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionAccepted)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "session_event.SessionAcceptFailed", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionAcceptFailed)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "session_event.SessionConnected", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionConnected)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "session_event.SessionConnectFailed", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionConnectFailed)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "session_event.SessionError", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionError)
        fmt.Println(msg)
    })

    socket.RegisterMessage(server, "session_event.SessionClosed", func(content interface{}, ses socket.Session) {
        msg := content.(sessproto.SessionClosed)
        fmt.Println(msg)
    })

    queue.StartLoop()
}

func sessionprotoRegisterMessage() {
    fmt.Println("session_event RegisterMessage")
    // session_event.proto
    socket.RegisterMessageMeta("sessevent.SessionAccepted", (*sessproto.SessionAccepted)(nil), 2136350511)
    socket.RegisterMessageMeta("sessevent.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1213847952)
    socket.RegisterMessageMeta("sessevent.SessionConnected", (*sessproto.SessionConnected)(nil), 4228538224)
    socket.RegisterMessageMeta("sessevent.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1278926828)
    socket.RegisterMessageMeta("sessevent.SessionClosed", (*sessproto.SessionClosed)(nil), 2830250790)
    socket.RegisterMessageMeta("sessevent.SessionError", (*sessproto.SessionError)(nil), 3227768243)

}

func loginprotoRegisterMessage()  {
    fmt.Println("login_event RegisterMessage")
    // login_event.proto
    socket.RegisterMessageMeta("login_event.PLoginByPassport", (*loginproto.PLoginByPassport)(nil), 3176521479)
    socket.RegisterMessageMeta("login_event.PLoginByPassportRes", (*loginproto.PLoginByPassportRes)(nil), 1894287111)
    socket.RegisterMessageMeta("login_event.PLoginByToken", (*loginproto.PLoginByToken)(nil), 3676221678)
    socket.RegisterMessageMeta("login_event.PLoginByTokenRes", (*loginproto.PLoginByTokenRes)(nil), 2978710459)
    socket.RegisterMessageMeta("login_event.PLoginLogout", (*loginproto.PLoginLogout)(nil), 182945212)
    socket.RegisterMessageMeta("login_event.PLoginLogoutRes", (*loginproto.PLoginLogoutRes)(nil), 3420227263)
    socket.RegisterMessageMeta("login_event.PLoginPing", (*loginproto.PLoginPing)(nil), 1951739067)
    socket.RegisterMessageMeta("login_event.PLoginPingRes", (*loginproto.PLoginPingRes)(nil), 3845948673)
    socket.RegisterMessageMeta("login_event.PLoginKickOff", (*loginproto.PLoginKickOff)(nil), 4292429949)
}