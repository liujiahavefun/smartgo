package main

import (
    sessproto "smartgo/proto/session_event"
    loginproto "smartgo/proto/login_event"
    "smartgo/libs/socket"
    "fmt"
)

var (
    gServer     socket.Server
    gEventQueue socket.EventQueue
)

func start(address string) {
    sessionprotoRegisterMessage()
    loginprotoRegisterMessage()
    socket.Init()

    runServer(address)
}

func sessionprotoRegisterMessage() {
    logInfo("Register Message for session_event")
    // session_event.proto
    socket.RegisterMessageMeta("session_event.SessionAccepted", (*sessproto.SessionAccepted)(nil), 348117910)
    socket.RegisterMessageMeta("session_event.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1978788392)
    socket.RegisterMessageMeta("session_event.SessionConnected", (*sessproto.SessionConnected)(nil), 3543838007)
    socket.RegisterMessageMeta("session_event.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1720533237)
    socket.RegisterMessageMeta("session_event.SessionClosed", (*sessproto.SessionClosed)(nil), 90181607)
    socket.RegisterMessageMeta("session_event.SessionError", (*sessproto.SessionError)(nil), 1937281175)

}

func loginprotoRegisterMessage()  {
    logInfo("Register Message for login_event")
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

func global_log(session socket.Session)  {
    logInfo(gSessionMgr, sessionInfo(session))
}

func sessionInfo(session socket.Session) string  {
    if session == nil {
        return fmt.Sprintf("[Session] empty")
    }
    return fmt.Sprintf("[Session], %v, %v", session.GetID(), session.FromPeer().Name())
}

func runServer(address string) {
    gEventQueue = socket.NewEventQueue()
    gServer = socket.NewTcpServer(gEventQueue).Start(address)

    //处理session消息
    socket.RegisterMessage(gServer, "session_event.SessionAccepted", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*sessproto.SessionAccepted)
        if !ok || msg == nil {
            logWarning("Server: recv invalid SessionAccepted message")
            return
        }
        handleSessionAccepted(msg, session)
    })

    socket.RegisterMessage(gServer, "session_event.SessionAcceptFailed", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*sessproto.SessionAcceptFailed)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid SessionAcceptFailed message")
            return
        }
        handleSessionAcceptFailed(msg, session)
    })

    socket.RegisterMessage(gServer, "session_event.SessionConnected", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*sessproto.SessionConnected)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid SessionConnected message")
            return
        }
        handleSessionConnected(msg, session)
    })

    socket.RegisterMessage(gServer, "session_event.SessionConnectFailed", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*sessproto.SessionConnectFailed)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid SessionConnectFailed message")
            return
        }
        handleSessionConnectFailed(msg, session)
    })

    socket.RegisterMessage(gServer, "session_event.SessionError", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*sessproto.SessionError)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid SessionError message")
            return
        }
        handleSessionError(msg, session)
    })

    socket.RegisterMessage(gServer, "session_event.SessionClosed", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*sessproto.SessionClosed)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid SessionClosed message")
            return
        }
        handleSessionClosed(msg, session)
    })

    //处理login消息
    socket.RegisterMessage(gServer, "login_event.PLoginByPassport", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*loginproto.PLoginByPassport)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid PLoginByPassport message")
            return
        }
        handleLoginByPassport(msg, session)
    })

    socket.RegisterMessage(gServer, "login_event.PLoginByToken", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*loginproto.PLoginByToken)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid PLoginByToken message")
            return
        }
        handleLoginByToken(msg, session)
    })

    socket.RegisterMessage(gServer, "login_event.PLoginLogout", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*loginproto.PLoginLogout)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid PLoginLogout message")
            return
        }
        handleLoginLogout(msg, session)
    })

    socket.RegisterMessage(gServer, "login_event.PLoginPing", func(content interface{}, session socket.Session) {
        global_log(session)
        msg, ok := content.(*loginproto.PLoginPing)
        if !ok || msg == nil  {
            logWarning("Server: recv invalid PLoginPing message")
            return
        }
        handleLoginPing(msg, session)
    })

    gEventQueue.StartLoop()
}