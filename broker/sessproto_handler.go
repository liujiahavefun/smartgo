package main

import (
    "time"

    sessproto "smartgo/proto/session_event"
    "smartgo/libs/socket"
)

func handleSessionAccepted(msg *sessproto.SessionAccepted, session socket.Session)  {
    logInfo("Server: recv SessionAccepted message, enter")
    logInfo("Server: recv SessionAccepted message, leave")
}

func handleSessionAcceptFailed(msg *sessproto.SessionAcceptFailed, session socket.Session)  {
    logInfo("Server: recv SessionAcceptFailed message, enter")
    logInfo("Server: recv SessionAcceptFailed message, leave")
}

func handleSessionConnected(msg *sessproto.SessionConnected, session socket.Session)  {
    logInfo("Server: recv SessionConnected message, enter")
    gSessionMgr.onSessionConnect(session)
    addLoginCheckTask(session)
    logInfo("Server: recv SessionConnected message, leave")
}

func handleSessionConnectFailed(msg *sessproto.SessionConnectFailed, session socket.Session)  {
    logInfo("Server: recv SessionConnectFailed message, enter")
    logInfo("Server: recv SessionConnectFailed message, leave")
}

func handleSessionError(msg *sessproto.SessionError, session socket.Session)  {
    logInfo("Server: recv SessionError message, enter")
    gSessionMgr.onSessionClose(session)
    logInfo("Server: recv SessionError message, leave")
}

func handleSessionClosed(msg *sessproto.SessionClosed, session socket.Session)  {
    logInfo("Server: recv SessionClosed message, enter")
    gSessionMgr.onSessionClose(session)
    logInfo("Server: recv SessionClosed message, leave")
}

func addLoginCheckTask(session socket.Session) {
    logInfo("Server: to add login check timer")
    socket.NewTimer(gEventQueue, time.Duration(gBrokerConfig.Timeout4Login)*time.Second, func(t *socket.Timer) {
        logInfo("login timer check")
        t.Stop()
        if logined, err := getSessionParamAsBool(session, SESSION_LOGINED); err != nil || logined == false {
            connected_time, _ := getSessionParamAsInt64(session, SESSION_LOGINED_TIME)
            logWarning("Stop session cause not logined after connected in 3 seconds, conneted time", connected_time)
            session.Close()
        }
    })
}