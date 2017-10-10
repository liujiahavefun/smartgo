package main

import (
    "time"

    sessproto "smartgo/proto/session_event"
    "smartgo/libs/socket"
)

func handleSessionAccepted(msg *sessproto.SessionAccepted, session socket.Session)  {
    logInfo("Server: recv SessionAccepted message")
}

func handleSessionAcceptFailed(msg *sessproto.SessionAcceptFailed, session socket.Session)  {
    logInfo("Server: recv SessionAcceptFailed message")
}

func handleSessionConnected(msg *sessproto.SessionConnected, session socket.Session)  {
    logInfo("Server: recv SessionConnected message")
    gSessionMgr.onSessionConnect(session)
    addLoginCheckTask(session)
}

func handleSessionConnectFailed(msg *sessproto.SessionConnectFailed, session socket.Session)  {
    logInfo("Server: recv SessionConnectFailed message")
}

func handleSessionError(msg *sessproto.SessionError, session socket.Session)  {
    logInfo("Server: recv SessionError message")
    gSessionMgr.onSessionClose(session)
}

func handleSessionClosed(msg *sessproto.SessionClosed, session socket.Session)  {
    logInfo("Server: recv SessionClosed message")
    gSessionMgr.onSessionClose(session)
}

func addLoginCheckTask(session socket.Session) {
    logInfo("Server: to add login check timer")
    socket.NewTimer(gEventQueue, 1*time.Second, func(t *socket.Timer) {
        logInfo("login timer check")
        t.Stop()
        if logined, err := getSessionParamAsBool(session, SESSION_LOGINED); err != nil || logined == false {
            connected_time, _ := getSessionParamAsInt64(session, SESSION_LOGINED_TIME)
            logWarning("Stop session cause not logined after connected in 3 seconds, conneted time", connected_time)
            session.Close()
        }
    })
}