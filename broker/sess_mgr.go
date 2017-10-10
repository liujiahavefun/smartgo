package main

import (
    "smartgo/libs/utils"
    "smartgo/libs/socket"
    "time"
    "fmt"
)

const (
    SESSION_CONNECTED      = "session_connected"
    SESSION_CONNECTED_TIME = "session_connected_time"
    SESSION_LOGINED        = "session_logined"
    SESSION_LOGINED_TIME   = "session_logined_time"
    SESSION_USER_ID        = "session_uid"
    SESSION_TOKEN          = "session_token"
    SESSION_PING_TIMER     = "session_ping_timer"
    SESSION_LAST_PING_TIME = "session_last_ping_time"
)

var (
    gSessionMgr *SessionMgr
)

type SessionMgr struct {
    sessions        *utils.ConcurrentMap
    connected        *utils.ConcurrentMap
}

func NewSessionMgr() *SessionMgr  {
    return &SessionMgr{
        sessions:    utils.NewConcurrentMap(),
        connected:    utils.NewConcurrentMap(),
    }
}

func (self * SessionMgr) onSessionConnect(session socket.Session)  {
    key := session.FromPeer().Name()
    if oldSession, exist := self.getConnectedSession(key); exist && oldSession != nil {
        //客户端是不是重复登录啊？删掉当前的
        //log something
        self.onSessionClose(oldSession)
    }

    session.SetParam(SESSION_CONNECTED, true)
    session.SetParam(SESSION_CONNECTED_TIME, utils.CurrentTimeMillSecond())
    self.connected.Put(key, session)
}

func (self * SessionMgr) onSessionLogin(session socket.Session, uid, token string)  {
    session.SetParam(SESSION_LOGINED, true)
    session.SetParam(SESSION_LOGINED_TIME, utils.CurrentTimeMillSecond())
    session.SetParam(SESSION_USER_ID, uid)
    session.SetParam(SESSION_TOKEN, token)
    addPingCheckTask(session)

    self.connected.Remove(session.FromPeer().Name())
    self.sessions.Put(uid, session)
}

func (self * SessionMgr) onSessionLogout(session socket.Session)  {

}

func (self * SessionMgr) onSessionKickOff(session socket.Session)  {

}

func (self * SessionMgr) onSessionClose(session socket.Session)  {
    //如果已经login了，要logout先
    if logined, err := getSessionParamAsBool(session, SESSION_LOGINED); err == nil && logined == true {
        if uid, err := getSessionParamAsString(session, SESSION_USER_ID); err == nil {
            if token, err := getSessionParamAsString(session, SESSION_TOKEN); err == nil {
                logout(uid, token)
                self.sessions.Remove(uid)
            }
        }
    }

    self.connected.Remove(session.FromPeer().Name())
    session.SetParam(SESSION_CONNECTED, false)
    session.SetParam(SESSION_LOGINED, false)

    if timer, err := getSessionParamAsTimer(session, SESSION_PING_TIMER); err == nil && timer != nil {
        timer.Stop()
    }

    //注意这里不要再Close了，收到此回调时已经close了，否则会无限循环
    //session.Close()
}

func (self * SessionMgr) String() (str string)  {
    return fmt.Sprintf("connected session %v, logined session %v", self.connected.Size(), self.sessions.Size())
}

func (self * SessionMgr) getConnectedSession(key string) (session socket.Session, exist bool)  {
    v, exist := self.connected.Get(key)
    if exist == false || v == nil {
        return nil, false
    }

    if session, ok := v.(socket.Session); !ok {
        //rarely case, however it's better to log something
        return nil, true
    }else {
        return session, true
    }
}

func (self * SessionMgr) getLoginedSession(key string) (session socket.Session, exist bool)  {
    v, exist := self.sessions.Get(key)
    if exist == false || v == nil {
        return nil, false
    }

    if session, ok := v.(socket.Session); !ok {
        //rarely case, however it's better to log something
        return nil, true
    }else {
        return session, true
    }
}

func addPingCheckTask(session socket.Session) {
    logInfo("Server: to add ping check timer")
    timer := socket.NewTimer(gEventQueue, 2*time.Second, func(t *socket.Timer) {
        logInfo("ping timer check")
        if last_ping_time, err := getSessionParamAsInt64(session, SESSION_LAST_PING_TIME);err != nil || utils.CurrentTimeMillSecond() - last_ping_time > 1*1000 {
            logWarning("Ping timeout, last ping is ", last_ping_time)
            session.Close()
        }
    })
    session.SetParam(SESSION_PING_TIMER, timer)
    session.SetParam(SESSION_LAST_PING_TIME, utils.CurrentTimeMillSecond())
}