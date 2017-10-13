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
    socket.RegisterMessageMeta("session_event.SessionAccepted", (*sessproto.SessionAccepted)(nil), uint32(PSessionAccepted_uri))
    socket.RegisterMessageMeta("session_event.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), uint32(PSessionAcceptedFailed_uri))
    socket.RegisterMessageMeta("session_event.SessionConnected", (*sessproto.SessionConnected)(nil), uint32(PSessionConnected_uri))
    socket.RegisterMessageMeta("session_event.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), uint32(PSessionConnectedFailed_uri))
    socket.RegisterMessageMeta("session_event.SessionClosed", (*sessproto.SessionClosed)(nil), uint32(PSessionClosed_uri))
    socket.RegisterMessageMeta("session_event.SessionError", (*sessproto.SessionError)(nil), uint32(PSessionError_uri))
}

func loginprotoRegisterMessage()  {
    logInfo("Register Message for login_event")
    // login_event.proto
    socket.RegisterMessageMeta("login_event.PLoginByPassport", (*loginproto.PLoginByPassport)(nil), uint32(PLoginByPassport_uri))
    socket.RegisterMessageMeta("login_event.PLoginByPassportRes", (*loginproto.PLoginByPassportRes)(nil), uint32(PLoginByPassportRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginByToken", (*loginproto.PLoginByToken)(nil), uint32(PLoginByUid_uri))
    socket.RegisterMessageMeta("login_event.PLoginByTokenRes", (*loginproto.PLoginByTokenRes)(nil), uint32(PLoginByUidRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginLogout", (*loginproto.PLoginLogout)(nil), uint32(PLoginLogout_uri))
    socket.RegisterMessageMeta("login_event.PLoginLogoutRes", (*loginproto.PLoginLogoutRes)(nil), uint32(PLoginLogoutRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginPing", (*loginproto.PLoginPing)(nil), uint32(PLoginPing_uri))
    socket.RegisterMessageMeta("login_event.PLoginPingRes", (*loginproto.PLoginPingRes)(nil), uint32(PLoginPingRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginKickOff", (*loginproto.PLoginKickOff)(nil), uint32(PLoginKickOff_uri))
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

    socket.RegisterDefault(gServer, func(msgId uint32, data []byte, session socket.Session) {
        logInfo("Server: recv unregistered message, id", msgId, "len", len(data))
    })

    gEventQueue.StartLoop()
}